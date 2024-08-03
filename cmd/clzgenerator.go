package cmd

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/kelveny/mockcompose/pkg/gogen"
	"github.com/kelveny/mockcompose/pkg/gosyntax"
	gosyntaxtyp "github.com/kelveny/mockcompose/pkg/gosyntax"

	"github.com/kelveny/mockcompose/pkg/gotype"
	"github.com/kelveny/mockcompose/pkg/logger"
)

const (
	MATCH_NONE matchType = iota
	MATCH_CLONE
	MATCH_MOCK
)

type matchType int

type classMethodGenerator struct {
	clzName     string // name of the class that implements class interface
	mockPkgName string // package name that mocking class resides
	mockName    string // the mocking composite class name

	methodsToClone []string // method function names that need to be cloned in mocking class
	methodsToMock  []string // method function names that need to be mocked
}

type generatorContext struct {
	mockedFunctions map[string]any

	// class type declaraion string -> method name -> *ReceiverSpec
	clzMethods map[string]map[string]*gosyntax.ReceiverSpec
}

func (c *generatorContext) hasFunctionMocked(fnName string) bool {
	if len(c.mockedFunctions) > 0 {
		if _, ok := c.mockedFunctions[fnName]; ok {
			return true
		}
	}
	return false
}

func (c *generatorContext) recordMockedFunction(fnName string) {
	if c.mockedFunctions == nil {
		c.mockedFunctions = make(map[string]any)
	}
	c.mockedFunctions[fnName] = struct{}{}
}

//go:generate mockcompose -n gctx_findClassMethods -c generatorContext -real findClassMethods,gosyntax
func (c *generatorContext) findClassMethods(
	clzTypeDeclString string,
	fset *token.FileSet,
	f *ast.File,
) map[string]*gosyntaxtyp.ReceiverSpec {
	if c.clzMethods == nil {
		c.clzMethods = make(map[string]map[string]*gosyntaxtyp.ReceiverSpec)
	}

	if _, ok := c.clzMethods[clzTypeDeclString]; !ok {
		c.clzMethods[clzTypeDeclString] = gosyntax.FindClassMethods(clzTypeDeclString, fset, f)
	}

	return c.clzMethods[clzTypeDeclString]
}

// use compiler to enforce interface compliance
var _ parsedFileGenerator = (*classMethodGenerator)(nil)

// match checks if a FuncDecl matches condition
func (g *classMethodGenerator) match(fnSpec *ast.FuncDecl) (bool, matchType) {
	if fnSpec.Recv != nil {
		t := fnSpec.Recv.List[0].Type

		if expr, ok := t.(*ast.StarExpr); ok {
			recvrClzName := expr.X.(*ast.Ident).Name
			if g.clzName == recvrClzName {
				if matchType := g.matchNameInConfig(fnSpec.Name.Name); matchType != MATCH_NONE {
					return true, matchType
				}
			}
		} else {
			recvrClzName := t.(*ast.Ident).Name
			if g.clzName == recvrClzName {
				if matchType := g.matchNameInConfig(fnSpec.Name.Name); matchType != MATCH_NONE {
					return true, matchType
				}
			}
		}
	} else {
		matchType := g.matchNameInConfig(fnSpec.Name.Name)
		if matchType != MATCH_NONE {
			return true, matchType
		}
	}

	return false, MATCH_NONE
}

func getReceiverTypeName(fnSpec *ast.FuncDecl) string {
	t := fnSpec.Recv.List[0].Type

	if expr, ok := t.(*ast.StarExpr); ok {
		return expr.X.(*ast.Ident).Name
	} else {
		return t.(*ast.Ident).Name
	}
}

func changeReceiverTypeName(fnSpec *ast.FuncDecl, name string) {
	t := fnSpec.Recv.List[0].Type

	if expr, ok := t.(*ast.StarExpr); ok {
		expr.X.(*ast.Ident).Name = name
	} else {
		t.(*ast.Ident).Name = name
	}
}

func (g *classMethodGenerator) matchNameInConfig(fnName string) matchType {
	if len(g.methodsToClone) > 0 {
		for _, name := range g.methodsToClone {
			// in format of methodName,pkg1=mockPkg1:pkg2=mockPkg2
			name = strings.Split(name, ",")[0]
			if name == fnName {
				return MATCH_CLONE
			}
		}
	}

	if len(g.methodsToMock) > 0 {
		for _, name := range g.methodsToMock {
			if name == fnName {
				return MATCH_MOCK
			}
		}
	}

	return MATCH_NONE
}

func (g *classMethodGenerator) getMethodOverrides(
	callerPkg string,
	fnName string,
	calleeVisitor *gosyntax.CalleeVisitor,
	receiver string,
) map[string]string {
	if calleeVisitor.ReceiverName() != "" {
		receiver = calleeVisitor.ReceiverName()
	}

	for _, name := range g.methodsToClone {
		// in format of methodName,(pkg1[]=mockPkg1]:pkg2[=mockPkg2],
		tokens := strings.Split(name, ",")
		if tokens[0] == fnName {
			if len(tokens) > 1 {
				pairs := strings.Split(tokens[1], ":")

				overrides := make(map[string]string)
				for _, pair := range pairs {
					kv := strings.Split(pair, "=")
					if len(kv) == 2 {
						if kv[0] != "." && kv[0] != "this" {
							overrides[kv[0]] = kv[1]
						} else {
							logger.Log(logger.ERROR, "invalid package override usage: %s", pair)
						}
					} else {
						switch kv[0] {
						case ".":
							// override all callee functions within the same package
							for _, calleeFn := range calleeVisitor.GetThisPackageCallees() {
								overrides[calleeFn] = fmt.Sprintf(
									"%s.%s.%s",
									receiver,
									g.getMockedPackageClzName(callerPkg, kv[0], fnName),
									calleeFn,
								)
							}
						case "this":
							// for peer method callee, we don't need extra override
						default:
							// other package callees, override with a mocked package class
							overrides[kv[0]] = fmt.Sprintf(
								"&%s.%s",
								receiver,
								g.getMockedPackageClzName(callerPkg, kv[0], fnName),
							)
						}
					}
				}
				return overrides
			}
		}
	}
	return nil
}

func (g *classMethodGenerator) getAutoMockCalleeConfig(
	fnName string,
) (autoMockPeers bool, autoMockPkgs []string) {
	for _, name := range g.methodsToClone {
		// in format of methodName,(pkg1[]=mockPkg1]:pkg2[=mockPkg2],
		tokens := strings.Split(name, ",")
		if tokens[0] == fnName {
			if len(tokens) > 1 {
				pairs := strings.Split(tokens[1], ":")

				for _, pair := range pairs {
					kv := strings.Split(pair, "=")
					if len(kv) == 1 {
						switch kv[0] {
						case "this":
							autoMockPeers = true
						default:
							autoMockPkgs = append(autoMockPkgs, kv[0])
						}
					}
				}
			}
		}
	}
	return
}

func (g *classMethodGenerator) composeMock(
	generatorCtx *generatorContext,
	writer io.Writer,
	fset *token.FileSet,
	fnSpec *ast.FuncDecl,
) {
	if !generatorCtx.hasFunctionMocked(fnSpec.Name.Name) {
		gogen.MockFunc(
			writer,
			g.mockPkgName,
			g.mockName,
			fset,
			fnSpec.Name.Name,
			fnSpec.Type.Params,
			fnSpec.Type.Results,
			nil,
		)

		generatorCtx.recordMockedFunction(fnSpec.Name.Name)
	}
}

func (g *classMethodGenerator) generate(
	writer io.Writer,
	file *ast.File,
) error {
	var buf bytes.Buffer

	fset := token.NewFileSet()
	if ok, autoMockPkgs := g.generateInternal(&buf, fset, file); ok {
		// reload generated content to process generated code the second time
		f, err := parser.ParseFile(fset, "", buf.Bytes(), parser.ParseComments)
		if err != nil {
			logger.Log(logger.ERROR, "Internal error: %s\n\n%s\n", err, buf.String())
			return err
		}

		// remove unused imports
		var cleanedImports []gosyntax.ImportSpec = []gosyntax.ImportSpec{
			{
				Name: "mock",
				Path: "github.com/stretchr/testify/mock",
			},
		}
		cleanedImports = gogen.CleanImports(f, cleanedImports)

		// compose final output
		fmt.Fprintf(writer, header, g.mockPkgName)

		gogen.WriteImportDecls(writer, cleanedImports)
		fmt.Fprintf(writer, compositeClzTemplateBegin, g.mockName, g.clzName, "mock.Mock")
		for _, mockedPkgClz := range autoMockPkgs {
			fmt.Fprintf(writer, "	%s\n", mockedPkgClz)
		}
		fmt.Fprint(writer, compositeClzTemplateEnd)

		for _, mockedPkgClz := range autoMockPkgs {
			fmt.Fprintf(writer, mockClzTemplate, mockedPkgClz, "mock.Mock")
		}

		gogen.WriteFuncDecls(writer, fset, f)
	}

	return nil
}

func (g *classMethodGenerator) generateInternal(
	writer io.Writer,
	fset *token.FileSet,
	file *ast.File,
) (generated bool, autoMockPkgs []string) {
	writer.Write([]byte(fmt.Sprintf("package %s\n\n", g.mockPkgName)))

	generatorCtx := &generatorContext{}
	imports := gosyntax.GetFileImportsAsMap(file)

	if len(file.Decls) > 0 {
		for _, d := range file.Decls {
			if fnSpec, ok := d.(*ast.FuncDecl); ok {
				matched, matchType := g.match(fnSpec)
				if matched {
					generated = true
				}

				if matchType == MATCH_CLONE {
					// check if we need to clone a method function or a ordinary function
					receiverSpec := gosyntax.FuncDeclReceiverSpec(fset, fnSpec)
					if receiverSpec != nil {
						//
						// clone a method function
						//

						// find out callee situation
						clzMethods := generatorCtx.findClassMethods(receiverSpec.TypeDecl, fset, file)
						v := gosyntax.NewCalleeVisitor(
							imports,
							clzMethods,
							receiverSpec.Name,
							fnSpec.Name.Name,
						)
						ast.Walk(v, fnSpec.Body)
						v.SanitizeCallees(imports)

						overrides := g.getMethodOverrides(file.Name.Name, fnSpec.Name.Name, v, "")

						n := getReceiverTypeName(fnSpec)
						changeReceiverTypeName(fnSpec, g.mockName)
						gogen.WriteFuncWithLocalOverrides(
							writer,
							fset,
							fnSpec,
							"",
							fnSpec.Name.Name,
							overrides,
						)
						changeReceiverTypeName(fnSpec, n)

						// pkgs will be in order of how it is defined in "-real,<pkg1>:<pkg2>""
						autoMockPeer, pkgs := g.getAutoMockCalleeConfig(fnSpec.Name.Name)
						if autoMockPeer {
							// peer callee in order of how it is declared in file
							g.generateMethodPeerCallees(generatorCtx, writer, fset, file, fnSpec, v)
						}

						if len(pkgs) > 0 {
							// package will be mocked as a struct type, mockedPkgs contains the names of these mocked structs
							mockedPkgs := g.generateFuncCallees(writer, fset, file, fnSpec, v, pkgs)
							if len(mockedPkgs) > 0 {
								autoMockPkgs = append(autoMockPkgs, mockedPkgs...)
							}
						}
					} else {
						//
						// clone a matched function
						//
						v := gosyntax.NewCalleeVisitor(
							imports,
							nil,
							"",
							fnSpec.Name.Name,
						)
						ast.Walk(v, fnSpec.Body)
						v.SanitizeCallees(imports)

						overrides := g.getMethodOverrides(file.Name.Name, fnSpec.Name.Name, v, "m")

						// create an artificial receiver
						gogen.WriteFuncWithLocalOverrides(
							writer,
							fset,
							fnSpec,
							fmt.Sprintf("(m *%s)", g.mockName),
							fnSpec.Name.Name,
							overrides,
						)

						// pkgs will be in order of how it is defined in "-real,<pkg1>:<pkg2>""
						_, pkgs := g.getAutoMockCalleeConfig(fnSpec.Name.Name)

						if len(pkgs) > 0 {
							// package will be mocked as a struct type, mockedPkgs contains the names of these mocked structs
							mockedPkgs := g.generateFuncCallees(writer, fset, file, fnSpec, v, pkgs)
							if len(mockedPkgs) > 0 {
								autoMockPkgs = append(autoMockPkgs, mockedPkgs...)
							}
						}
					}
				} else if matchType == MATCH_MOCK {
					// generate mocked method
					g.composeMock(generatorCtx, writer, fset, fnSpec)
				}
			} else {
				// for any non-function declaration, export only imports
				if dd, ok := d.(*ast.GenDecl); ok && dd.Tok == token.IMPORT {
					format.Node(writer, fset, d)
					writer.Write([]byte("\n\n"))
				}
			}
		}
	}

	return
}

func (g *classMethodGenerator) generateMethodPeerCallees(
	generatorCtx *generatorContext,
	writer io.Writer,
	fset *token.FileSet,
	file *ast.File,
	callerFnSpec *ast.FuncDecl,
	calleeVisitor *gosyntax.CalleeVisitor,
) {
	if len(calleeVisitor.GetPeerCallees()) > 0 {
		for _, peerMethod := range calleeVisitor.GetPeerCallees() {

			// if peer method is not in explicitly specified mocking configuration,
			// generate it automatically
			if !slices.Contains(g.methodsToMock, peerMethod) {
				gosyntax.ForEachFuncDeclInFile(file, func(fnSpec *ast.FuncDecl) {
					if fnSpec.Name.Name == peerMethod &&
						gosyntax.ReceiverDeclString(fset, callerFnSpec.Recv) == gosyntax.ReceiverDeclString(fset, fnSpec.Recv) {
						g.composeMock(generatorCtx, writer, fset, fnSpec)
					}
				})
			}
		}
	}
}

func (g *classMethodGenerator) generateFuncCallees(
	writer io.Writer,
	_ *token.FileSet,
	file *ast.File,
	callerFnSpec *ast.FuncDecl,
	calleeVisitor *gosyntax.CalleeVisitor,
	pkgs []string,
) []string {
	mockedPkgs := []string{}

	imports := gosyntax.GetFileImportsAsMap(file)

	for _, pkg := range pkgs {
		mockedPkg := g.getMockedPackageClzName(file.Name.Name, pkg, callerFnSpec.Name.Name)

		var callees []string
		if pkg == "." {
			callees = calleeVisitor.GetThisPackageCallees()
		} else {
			callees = calleeVisitor.GetOtherPackageCallees()[pkg]
		}

		for _, callee := range callees {
			calleeSpec, err := gotype.GetFuncTypeSpec(imports[pkg], callee, g.mockPkgName)
			if err == nil {
				gogen.GenerateFuncMock(
					writer,
					g.mockPkgName,
					mockedPkg,
					callee,
					calleeSpec.FieldInfo,
					calleeSpec.ReturnInfo,
					calleeSpec.Signature,
				)
			}
		}

		mockedPkgs = append(mockedPkgs, mockedPkg)
	}
	return mockedPkgs
}

func (g *classMethodGenerator) getMockedPackageClzName(
	callerPkg string,
	pkgNameToMock string,
	funcName string,
) string {
	// scope mocked package class name at per-caller per-ref-package basis
	if pkgNameToMock != "." {
		return fmt.Sprintf("mock_%s_%s_%s", g.mockName, funcName, pkgNameToMock)
	}

	return fmt.Sprintf("mock_%s_%s_%s", g.mockName, funcName, callerPkg)
}

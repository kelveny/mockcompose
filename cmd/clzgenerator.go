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

// use compiler to enforce interface compliance
var _ parsedFileGenerator = (*classMethodGenerator)(nil)

// match checks if a FuncDecl matches condition
func (g *classMethodGenerator) match(fnSpec *ast.FuncDecl) (bool, matchType) {
	if fnSpec.Recv != nil {
		t := fnSpec.Recv.List[0].Type

		if expr, ok := t.(*ast.StarExpr); ok {
			recvrClzName := expr.X.(*ast.Ident).Name
			if g.clzName == recvrClzName {
				if matchType := g.matchMethod(fnSpec.Name.Name); matchType != MATCH_NONE {
					return true, matchType
				}
			}
		} else {
			recvrClzName := t.(*ast.Ident).Name
			if g.clzName == recvrClzName {
				if matchType := g.matchMethod(fnSpec.Name.Name); matchType != MATCH_NONE {
					return true, matchType
				}
			}
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

func (g *classMethodGenerator) matchMethod(fnName string) matchType {
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
	fnName string,
	calleeVisitor *gosyntax.MethodCalleeVisitor,
) map[string]string {
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
									calleeVisitor.ReceiverName(),
									g.getMockedPackageClzName(kv[0], fnName),
									calleeFn,
								)
							}
						case "this":
							// for peer method callee, we don't need extra override
						default:
							// other package callees, override with a mocked package class
							overrides[kv[0]] = fmt.Sprintf(
								"%s.%s",
								calleeVisitor.ReceiverName(),
								g.getMockedPackageClzName(kv[0], fnName),
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

func (g *classMethodGenerator) getMethodAutoMockCalleeConfig(
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
	writer io.Writer,
	fset *token.FileSet,
	fnSpec *ast.FuncDecl,
) {
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
}

func (g *classMethodGenerator) generate(
	writer io.Writer,
	file *ast.File,
) error {

	var buf bytes.Buffer

	fset := token.NewFileSet()
	if ok, _ := g.generateInternal(&buf, fset, file); ok {
		// reload generated content to process generated code the second time
		f, err := parser.ParseFile(fset, "", buf.Bytes(), parser.ParseComments)
		if err != nil {
			logger.Log(logger.ERROR, "Internal error: %s\n\n%s\n", err, buf.String())
			return err
		}

		// remove unused imports
		var cleanedImports []gogen.ImportSpec = []gogen.ImportSpec{
			{
				Name: "mock",
				Path: "github.com/stretchr/testify/mock",
			},
		}
		cleanedImports = gogen.CleanImports(f, cleanedImports)

		// compose final output
		fmt.Fprintf(writer, header, g.mockPkgName)

		gogen.WriteImportDecls(writer, cleanedImports)
		fmt.Fprintf(writer, compositeClzTemplate, g.mockName, g.clzName, "mock.Mock")

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

	if len(file.Decls) > 0 {
		for _, d := range file.Decls {
			if fnSpec, ok := d.(*ast.FuncDecl); ok {
				matched, matchType := g.match(fnSpec)
				if matched {
					generated = true
				}

				if matchType == MATCH_CLONE {
					// clone receiver-modified method

					// find out callee situation
					receiverSpec := gosyntax.FuncDeclReceiverSpec(fset, fnSpec)
					clzMethods := gosyntax.FindClassMethods(receiverSpec.TypeDecl, fset, file)
					v := gosyntax.NewMethodCalleeVisitor(clzMethods, receiverSpec.Name, fnSpec.Name.Name)
					ast.Walk(v, fnSpec.Body)

					overrides := g.getMethodOverrides(fnSpec.Name.Name, v)

					n := getReceiverTypeName(fnSpec)
					changeReceiverTypeName(fnSpec, g.mockName)
					gogen.WriteFuncWithLocalOverrides(
						writer,
						fset,
						fnSpec,
						fnSpec.Name.Name,
						overrides,
					)
					changeReceiverTypeName(fnSpec, n)

					autoMockPeer, pkgs := g.getMethodAutoMockCalleeConfig(fnSpec.Name.Name)
					if autoMockPeer {
						g.generateMethodPeerCallees(writer, fset, file, fnSpec, v)
					}

					if len(pkgs) > 0 {
						pkgs := g.generateMethodFuncCallees(writer, fset, file, fnSpec, v)
						if len(pkgs) > 0 {
							autoMockPkgs = append(autoMockPkgs, pkgs...)
						}
					}

				} else if matchType == MATCH_MOCK {
					// generate mocked method
					g.composeMock(writer, fset, fnSpec)
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
	writer io.Writer,
	fset *token.FileSet,
	file *ast.File,
	callerFnSpec *ast.FuncDecl,
	calleeVisitor *gosyntax.MethodCalleeVisitor,
) {
	if len(calleeVisitor.GetPeerCallees()) > 0 {
		for _, peerMethod := range calleeVisitor.GetPeerCallees() {

			// if peer method is not in explicitly specified mocking configuration,
			// generate it automatically
			if !slices.Contains(g.methodsToMock, peerMethod) {

				gosyntax.ForEachFuncDeclInFile(file, func(fnSpec *ast.FuncDecl) {
					if fnSpec.Name.Name == peerMethod &&
						gosyntax.ReceiverDeclString(fset, callerFnSpec.Recv) == gosyntax.ReceiverDeclString(fset, fnSpec.Recv) {
						g.composeMock(writer, fset, fnSpec)
					}
				})
			}
		}
	}
}

func (g *classMethodGenerator) generateMethodFuncCallees(
	writer io.Writer,
	fset *token.FileSet,
	file *ast.File,
	callerFnSpec *ast.FuncDecl,
	calleeVisitor *gosyntax.MethodCalleeVisitor,
) []string {
	// ???
	return nil
}

func (g *classMethodGenerator) getMockedPackageClzName(
	pkgNameToMock string,
	funcName string,
) string {
	// scope mocked package class name at per-method per-package basis
	if pkgNameToMock != "." {
		return fmt.Sprintf("mock_%s_%s_%s", g.mockName, funcName, pkgNameToMock)
	}

	return fmt.Sprintf("mock_%s_%s_funcs", g.mockName, funcName)
}

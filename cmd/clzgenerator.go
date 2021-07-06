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

	"github.com/kelveny/mockcompose/pkg/gogen"
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

// match checks if a FuncDecl matches condition, if it is matched,
// change and prepare receiver type to the mocking target class
func (g *classMethodGenerator) match(fnSpec *ast.FuncDecl) (bool, matchType) {
	if fnSpec.Recv != nil {
		t := fnSpec.Recv.List[0].Type

		if expr, ok := t.(*ast.StarExpr); ok {
			recvrClzName := expr.X.(*ast.Ident).Name
			if g.clzName == recvrClzName {
				if matchType := g.matchMethod(fnSpec.Name.Name); matchType != MATCH_NONE {
					expr.X.(*ast.Ident).Name = g.mockName
					return true, matchType
				}
			}
		} else {
			recvrClzName := t.(*ast.Ident).Name
			if g.clzName == recvrClzName {
				if matchType := g.matchMethod(fnSpec.Name.Name); matchType != MATCH_NONE {
					t.(*ast.Ident).Name = g.mockName
					return true, matchType
				}
			}
		}
	}

	return false, MATCH_NONE
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

func (g *classMethodGenerator) getMethodOverrides(fnName string) map[string]string {
	for _, name := range g.methodsToClone {
		// in format of methodName,pkg1=mockPkg1:pkg2=mockPkg2
		tokens := strings.Split(name, ",")
		if tokens[0] == fnName {
			if len(tokens) > 1 {
				pairs := strings.Split(tokens[1], ":")

				overrides := make(map[string]string)
				for _, pair := range pairs {
					kv := strings.Split(pair, "=")
					if len(kv) == 2 {
						overrides[kv[0]] = kv[1]
					} else {
						logger.Log(logger.ERROR, "invalid configuration: -real %s\n", name)
					}
				}
				return overrides
			}
		}
	}
	return nil
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
	if g.generateInternal(&buf, fset, file) {
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
) bool {
	found := false

	writer.Write([]byte(fmt.Sprintf("package %s\n\n", g.mockPkgName)))

	if len(file.Decls) > 0 {
		for _, d := range file.Decls {
			if fnSpec, ok := d.(*ast.FuncDecl); ok {
				matched, matchType := g.match(fnSpec)
				if matched {
					found = true
				}

				if matchType == MATCH_CLONE {
					// clone receiver-modified method
					overrides := g.getMethodOverrides(fnSpec.Name.Name)
					gogen.WriteFuncWithLocalOverrides(
						writer,
						fset,
						fnSpec,
						fnSpec.Name.Name,
						overrides,
					)
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

	return found
}

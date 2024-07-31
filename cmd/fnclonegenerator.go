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
	"github.com/kelveny/mockcompose/pkg/gosyntax"
	"github.com/kelveny/mockcompose/pkg/logger"
)

type functionCloneGenerator struct {
	mockPkgName    string   // package name that cloned functions reside
	mockName       string   // name used to form generated file name
	methodsToClone []string // function names that need to be cloned
}

// use compiler to enforce interface compliance
var _ parsedFileGenerator = (*functionCloneGenerator)(nil)

// match checks if a FuncDecl matches condition, if matched,
// modify function name
func (g *functionCloneGenerator) match(fnSpec *ast.FuncDecl) bool {
	if fnSpec.Recv == nil {
		return g.matchMethod(fnSpec.Name.Name)
	}

	return false
}

func (g *functionCloneGenerator) matchMethod(fnName string) bool {
	if len(g.methodsToClone) > 0 {
		for _, name := range g.methodsToClone {
			// in format of methodName,pkg1=mockPkg1:pkg2=mockPkg2
			name = strings.Split(name, ",")[0]
			if name == fnName {
				return true
			}
		}
	}

	return false
}

func (g *functionCloneGenerator) getMethodOverrides(fnName string) map[string]string {
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

func (g *functionCloneGenerator) generate(
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
		var cleanedImports []gosyntax.ImportSpec = []gosyntax.ImportSpec{}
		cleanedImports = gogen.CleanImports(f, cleanedImports)

		// compose final output
		fmt.Fprintf(writer, header, g.mockPkgName)

		gogen.WriteImportDecls(writer, cleanedImports)
		gogen.WriteFuncDecls(writer, fset, f)
	}

	return nil
}

func (g *functionCloneGenerator) generateInternal(
	writer io.Writer,
	fset *token.FileSet,
	file *ast.File,
) bool {
	found := false

	writer.Write([]byte(fmt.Sprintf("package %s\n\n", g.mockPkgName)))

	if len(file.Decls) > 0 {
		for _, d := range file.Decls {
			if fnSpec, ok := d.(*ast.FuncDecl); ok {
				matched := g.match(fnSpec)
				if matched {
					found = true
				}

				if matched {
					overrides := g.getMethodOverrides(fnSpec.Name.Name)
					gogen.WriteFuncWithLocalOverrides(
						writer,
						fset,
						fnSpec,
						fnSpec.Name.Name+"_clone",
						overrides,
					)
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

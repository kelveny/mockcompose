package cmd

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"

	"github.com/kelveny/mockcompose/pkg/gogen"
	"github.com/kelveny/mockcompose/pkg/gosyntax"
	"github.com/kelveny/mockcompose/pkg/gotype"
	"github.com/kelveny/mockcompose/pkg/logger"
	"golang.org/x/tools/go/packages"
)

type functionMockGenerator struct {
	mockPkgName   string   // package name that mocking class resides
	mockName      string   // the mocking composite class name
	methodsToMock []string // function names that need to be mocked
	srcPkg        string
}

// use compiler to enforce interface compliance
var _ parsedFileGenerator = (*functionMockGenerator)(nil)
var _ loadedPackageGenerator = (*functionMockGenerator)(nil)

func (g *functionMockGenerator) generate(
	writer io.Writer,
	file *ast.File,
) error {
	var buf bytes.Buffer
	fset := token.NewFileSet()

	// first pass
	matchCount := 0
	var bufWriter io.Writer = &buf
	gosyntax.ForEachFuncDeclInFile(file, func(fnDecl *ast.FuncDecl) {
		if fnDecl.Recv == nil && g.match(fnDecl.Name.Name) {
			matchCount++
			if matchCount == 1 {
				bufWriter.Write([]byte(fmt.Sprintf("package %s\n\n", g.mockPkgName)))
			}

			gogen.MockFunc(
				bufWriter,
				g.mockPkgName,
				g.mockName,
				fset,
				fnDecl.Name.Name,
				fnDecl.Type.Params,
				fnDecl.Type.Results,
				nil,
			)
		}
	})

	// second pass
	if matchCount > 0 {
		// reload generated content to process generated code the second time
		f, err := parser.ParseFile(fset, "", buf.Bytes(), parser.ParseComments)
		if err != nil {
			logger.Log(logger.PROMPT, "Internal error: %s\n\n%s\n", err, buf.String())
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
		fmt.Fprintf(writer, mockClzTemplate, g.mockName, "mock.Mock")

		gogen.WriteFuncDecls(writer, fset, f)
	}

	return nil
}

func (g *functionMockGenerator) generateViaLoadedPackage(
	writer io.Writer,
	pkg *packages.Package,
) error {
	var buf bytes.Buffer
	fset := token.NewFileSet()

	// first pass
	matchCount := 0
	var bufWriter io.Writer = &buf

	gosyntax.ForEachFuncDeclInPackage(pkg, func(fnDecl *ast.FuncDecl) {
		if g.match(fnDecl.Name.Name) {
			matchCount++
			if matchCount == 1 {
				bufWriter.Write([]byte(fmt.Sprintf("package %s\n\n", g.mockPkgName)))

				imports := gogen.GetPackageImports(pkg)
				if g.srcPkg != "" {
					imports = append(imports, gosyntax.ImportSpec{
						Name: "",
						Path: g.srcPkg,
					})
				}
				gogen.WriteImportDecls(bufWriter, imports)
			}

			gogen.MockFunc(
				bufWriter,
				g.mockPkgName,
				g.mockName,
				fset,
				fnDecl.Name.Name,
				fnDecl.Type.Params,
				fnDecl.Type.Results,
				gotype.FindFuncSignature(pkg, fnDecl.Name.Name),
			)
		}
	})

	// second pass
	if matchCount > 0 {
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
		fmt.Fprintf(writer, mockClzTemplate, g.mockName, "mock.Mock")

		gogen.WriteFuncDecls(writer, fset, f)
	}

	return nil
}

func (g *functionMockGenerator) match(name string) bool {
	for _, n := range g.methodsToMock {
		if n == name {
			return true
		}
	}

	return false
}

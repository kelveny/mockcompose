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
	"github.com/kelveny/mockcompose/pkg/logger"
	"golang.org/x/tools/go/packages"
)

type interfaceMockGenerator struct {
	mockPkgName string // package name that mocking class resides
	mockName    string // the mocking composite class name
	intfName    string // interface name
}

// use compiler to enforce interface compliance
var _ parsedFileGenerator = (*interfaceMockGenerator)(nil)
var _ loadedPackageGenerator = (*interfaceMockGenerator)(nil)

func (g *interfaceMockGenerator) generate(
	writer io.Writer,
	file *ast.File,
) error {

	gosyntax.ForEachInterfaceDeclInFile(file,
		func(name string, methods []*ast.Field) {
			if name == g.intfName {
				imports := gogen.GetFileImports(file)
				g.generateInterfaceMock(writer, imports, methods)
			}
		},
	)
	return nil
}

func (g *interfaceMockGenerator) generateViaLoadedPackage(
	writer io.Writer,
	pkg *packages.Package,
) error {
	gosyntax.ForEachInterfaceDeclInPackage(pkg,
		func(name string, methods []*ast.Field) {
			if name == g.intfName {
				imports := gogen.GetPackageImports(pkg)
				g.generateInterfaceMock(writer, imports, methods)
			}
		},
	)
	return nil
}

func (g *interfaceMockGenerator) generateInterfaceMock(
	writer io.Writer,
	imports []gogen.ImportSpec,
	methods []*ast.Field,
) error {
	var buf bytes.Buffer

	fset := token.NewFileSet()
	if g.generateInterfaceMockInternal(&buf, fset, imports, methods) {
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
		fmt.Fprintf(writer, mockClzTemplate, g.mockName, "mock.Mock")

		gogen.WriteFuncDecls(writer, fset, f)
	}

	return nil
}

func (g *interfaceMockGenerator) generateInterfaceMockInternal(
	writer io.Writer,
	fset *token.FileSet,
	imports []gogen.ImportSpec,
	methods []*ast.Field,
) bool {
	writer.Write([]byte(fmt.Sprintf("package %s\n\n", g.mockPkgName)))

	gogen.WriteImportDecls(writer, imports)

	for _, method := range methods {
		if ftype, ok := method.Type.(*ast.FuncType); ok {
			gosyntax.MockFunc(
				writer,
				g.mockName,
				fset,
				method.Names[0].Name,
				ftype.Params,
				ftype.Results,
			)
		}
	}

	return true
}

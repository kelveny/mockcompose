package cmd

import (
	"go/ast"
	"io"

	"golang.org/x/tools/go/packages"
)

const (
	header = `//
// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
//
package %s
`
	compositeClzTemplate = `type %s struct {
	%s
	%s
}

`
	mockClzTemplate = `type %s struct {
	%s
}

`
)

type commandOptions struct {
	mockName *string
	mockPkg  *string
	clzName  *string
	intfName *string
	srcPkg   *string
	testOnly *bool

	methodsToClone stringSlice
	methodsToMock  stringSlice
}

type parsedFileGenerator interface {
	generate(writer io.Writer, file *ast.File) error
}

type loadedPackageGenerator interface {
	generateViaLoadedPackage(writer io.Writer, pkg *packages.Package) error
}

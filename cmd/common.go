package cmd

import (
	"go/ast"
	"io"

	"golang.org/x/tools/go/packages"
)

const (
	header = `// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
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

// must be public for it to be used in loading YAML configuration
type CommandOptions struct {
	MockName string `yaml:"name"`
	MockPkg  string `yaml:"mockPkg"`
	ClzName  string `yaml:"className"`
	IntfName string `yaml:"interfaceName"`
	SrcPkg   string `yaml:"sourcePkg"`
	TestOnly bool   `yaml:"testOnly"`

	// For example content: "functionThatUsesMultileGlobalFunctions,fmt=fmtMock:json=jsonMock",
	// it means to use "fmtMock" class as package "fmt", "jsonMock" class as package "json" in the cloned function
	//
	// For example content: "functionThatUsesMultileGlobalFunctions,this:.:fmt:json",
	// it means to mock peer callee methods (this as psudo package name), auto generated callee packages for "." package, "fmt" amd "json" package

	MethodsToClone []string `yaml:"real,flow"`

	MethodsToMock []string `yaml:"mock,flow"`
}

// must be public for it to be used in loading YAML configuration
type Config struct {
	Mockcompose []CommandOptions `yaml:"mockcompose,flow"`
}

type parsedFileGenerator interface {
	generate(writer io.Writer, file *ast.File) error
}

type loadedPackageGenerator interface {
	generateViaLoadedPackage(writer io.Writer, pkg *packages.Package) error
}

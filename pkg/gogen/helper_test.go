package gogen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

// for debugging purpose only
func TestTypesPackageFromPackagesPackage(t *testing.T) {
	assert := require.New(t)

	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax}

	// return packages.Package
	pkgs, err := packages.Load(cfg, "github.com/kelveny/mockcompose/test/libfn")
	if err != nil {
		t.FailNow()
	}

	assert.True(len(pkgs) == 1)

	for _, pkg := range pkgs {
		// pkg.Types points to types.Package
		assert.True(pkg.Types.Name() == "libfn")
		assert.True(pkg.Types.Path() == "github.com/kelveny/mockcompose/test/libfn")
	}
}

func TestGetPackageImports(t *testing.T) {
	assert := require.New(t)

	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, "github.com/kelveny/mockcompose/test/libfn")
	if err != nil {
		t.FailNow()
	}

	assert.True(len(pkgs) == 1)

	// return packages.Package
	imports := GetPackageImports(pkgs[0])

	assert.True(len(imports) > 0)

	for _, imp := range imports {
		fmt.Printf("package import, name: %s, path: %s\n", imp.Name, imp.Path)
	}
}

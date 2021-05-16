package gosyntax

import (
	"fmt"
	"go/ast"
	"testing"

	"golang.org/x/tools/go/packages"
)

// for debugging purpose only
func TestEachFuncInPackage(t *testing.T) {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax}

	pkgs, err := packages.Load(cfg, "filepath")
	if err != nil {
		return
	}

	for _, pkg := range pkgs {
		ForEachFuncDeclInPackage(pkg, func(fn *ast.FuncDecl) {
			fmt.Println(fn.Name.Name)
		})
	}
}

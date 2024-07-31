package gosyntax

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
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

func TestCalleeDetection(t *testing.T) {
	assert := require.New(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime.Caller failed")
	}
	cur, _ := filepath.Abs(filename)
	fooFile := filepath.Join(filepath.Dir(cur), "../../test/foo/foo.go")

	fset := token.NewFileSet()
	node, err := parser.ParseFile(
		fset,
		fooFile,
		nil,
		parser.ParseComments)

	imports := GetFileImportsAsMap(node)

	ForEachFuncDeclInFile(node, func(funcDecl *ast.FuncDecl) {
		receiverSpec := FuncDeclReceiverSpec(fset, funcDecl)

		if receiverSpec != nil && funcDecl.Name.Name == "Foo" && receiverSpec.TypeDecl == "*dummyFoo" {
			body := funcDecl.Body

			clzMethods := FindClassMethods("*dummyFoo", fset, node)

			v := NewCalleeVisitor(
				imports,
				clzMethods,
				receiverSpec.Name,  // receiver variable name
				funcDecl.Name.Name, // caller method function name
			)
			ast.Walk(v, body)

			assert.Equal([]string{"Bar"}, v.GetPeerCallees())
			assert.Equal([]string{"dummy"}, v.GetThisPackageCallees())
			assert.Equal(map[string][]string{
				"fmt": {"Printf"},
			}, v.GetOtherPackageCallees())
		}

		if receiverSpec == nil && funcDecl.Name.Name == "dummy" {
			body := funcDecl.Body

			v := NewCalleeVisitor(
				imports,
				nil,
				"",
				funcDecl.Name.Name, // caller function name
			)
			ast.Walk(v, body)

			assert.Equal(0, len(v.GetPeerCallees()))
			assert.Equal([]string{"fool"}, v.GetThisPackageCallees())
			assert.Equal(map[string][]string{
				"fmt": {"Print"},
			}, v.GetOtherPackageCallees())
		}
	})

	assert.NoError(err)
}

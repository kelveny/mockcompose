// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package cmd

import (
	"go/ast"
	"go/token"

	"github.com/kelveny/mockcompose/pkg/gosyntax"
	gosyntaxtyp "github.com/kelveny/mockcompose/pkg/gosyntax"
	"github.com/stretchr/testify/mock"
)

type gctx_findClassMethods struct {
	generatorContext
	mock.Mock
	mock_gctx_findClassMethods_findClassMethods_gosyntax
}

type mock_gctx_findClassMethods_findClassMethods_gosyntax struct {
	mock.Mock
}

func (c *gctx_findClassMethods) findClassMethods(clzTypeDeclString string, fset *token.FileSet, f *ast.File) map[string]*gosyntaxtyp.ReceiverSpec {
	gosyntax := &c.mock_gctx_findClassMethods_findClassMethods_gosyntax

	if c.clzMethods == nil {
		c.clzMethods = make(map[string]map[string]*gosyntaxtyp.ReceiverSpec)
	}
	if _, ok := c.clzMethods[clzTypeDeclString]; !ok {
		c.clzMethods[clzTypeDeclString] = gosyntax.FindClassMethods(clzTypeDeclString, fset, f)
	}
	return c.clzMethods[clzTypeDeclString]
}

func (m *mock_gctx_findClassMethods_findClassMethods_gosyntax) FindClassMethods(clzTypeDeclString string, fset *token.FileSet, f *ast.File) map[string]*gosyntax.ReceiverSpec {

	_mc_ret := m.Called(clzTypeDeclString, fset, f)

	var _r0 map[string]*gosyntax.ReceiverSpec

	if _rfn, ok := _mc_ret.Get(0).(func(string, *token.FileSet, *ast.File) map[string]*gosyntax.ReceiverSpec); ok {
		_r0 = _rfn(clzTypeDeclString, fset, f)
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(map[string]*gosyntax.ReceiverSpec)
		}
	}

	return _r0

}

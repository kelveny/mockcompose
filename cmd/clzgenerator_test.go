package cmd

import (
	"testing"

	"github.com/kelveny/mockcompose/pkg/gosyntax"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_generatorContext_findClassMethods_caching(t *testing.T) {
	assert := require.New(t)

	g := &gctx_findClassMethods{}

	g.mock_gctx_findClassMethods_findClassMethods_gosyntax.On(
		"FindClassMethods",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(
		map[string]*gosyntax.ReceiverSpec{
			"Foo": {
				Name:     "f",
				TypeDecl: "*foo",
			},
		},
	)

	// call it once
	methods := g.findClassMethods("*foo", nil, nil)
	assert.EqualValues(
		map[string]*gosyntax.ReceiverSpec{
			"Foo": {
				Name:     "f",
				TypeDecl: "*foo",
			},
		},
		methods,
	)

	// call it the second time
	methods = g.findClassMethods("*foo", nil, nil)
	assert.EqualValues(
		map[string]*gosyntax.ReceiverSpec{
			"Foo": {
				Name:     "f",
				TypeDecl: "*foo",
			},
		},
		methods,
	)

	// assert on caching behave
	g.mock_gctx_findClassMethods_findClassMethods_gosyntax.AssertNumberOfCalls(t, "FindClassMethods", 1)
}

func Test_generatorContext_findClassMethods_nil_return(t *testing.T) {
	assert := require.New(t)

	g := &gctx_findClassMethods{}

	g.mock_gctx_findClassMethods_findClassMethods_gosyntax.On(
		"FindClassMethods",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil)

	// call it once
	methods := g.findClassMethods("*foo", nil, nil)
	assert.Nil(methods)

	// call it the second time
	methods = g.findClassMethods("*foo", nil, nil)
	assert.Nil(methods)

	// assert on caching behave
	g.mock_gctx_findClassMethods_findClassMethods_gosyntax.AssertNumberOfCalls(t, "FindClassMethods", 1)
}

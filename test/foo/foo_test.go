package foo

import (
	"testing"

	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFoo(t *testing.T) {
	assert := require.New(t)

	fooObj := &testFoo{
		foo{
			name: "name of foo",
		},
		mock.Mock{},
	}

	fooObj.On("Bar").Return(false)

	s := fooObj.Foo()
	assert.True(s == "name of foo")
}

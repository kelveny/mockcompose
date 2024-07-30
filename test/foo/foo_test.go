package foo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFoo(t *testing.T) {
	assert := require.New(t)

	fooObj := &testFoo{
		foo: foo{
			name: "name of foo",
		},
	}

	// mock peer method bar() called from Foo()
	fooObj.On("Bar").Return(false)

	// mock a package level function
	fooObj.mock_testFoo_Foo_foo.On("dummy").Return()

	// mock a function from other package
	fooObj.mock_testFoo_Foo_fmt.On("Print", "Foo").Return(0, nil)

	// call real implementation code in Foo()
	s := fooObj.Foo()

	assert.True(s == "name of foo")
	fooObj.AssertNumberOfCalls(t, "Bar", 1)
	fooObj.mock_testFoo_Foo_foo.AssertNumberOfCalls(t, "dummy", 1)
	fooObj.mock_testFoo_Foo_fmt.AssertNumberOfCalls(t, "Print", 1)
}

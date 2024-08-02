package bar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFooBar(t *testing.T) {
	assert := require.New(t)

	fb := &fooBarMock{
		fooBar: fooBar{name: "TestFooBar"},
	}

	fb.On("order").Return(1).Once()
	fb.On("order").Return(2).Once()

	fb.mock_fooBarMock_BarFoo_bar.On("order").Return(2).Once()
	fb.mock_fooBarMock_BarFoo_bar.On("order").Return(1).Once()

	fb.On("Foo").Return("FooMocked")
	fb.On("Bar").Return("BarMocked")

	s1 := fb.FooBar()
	assert.Equal("TestFooBar: BarMockedFooMocked", s1)
	s2 := fb.BarFoo()
	assert.Equal(s1, s2)

	s1 = fb.FooBar()
	assert.Equal("TestFooBar: FooMockedBarMocked", s1)
	s2 = fb.BarFoo()
	assert.Equal(s1, s2)
}

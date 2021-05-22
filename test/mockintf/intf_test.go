package mockintf

import (
	"testing"

	"github.com/kelveny/mockcompose/test/foo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestVariadic(t *testing.T) {
	assert := require.New(t)

	si := SampleInterfaceImpl{}

	expected := "string1: value1, string2: value2, string3: value3"
	result := si.Variadic("string1: %s, string2: %s, string3: %s",
		"value1", "value2", "value3")

	assert.True(result == expected)
}

func TestVariadic2(t *testing.T) {
	assert := require.New(t)

	si := SampleInterfaceImpl{}

	expected := "string1: value1, integer2: 2, string3: value3"
	result := si.Variadic2("string1: %s, integer2: %d, string3: %s",
		"value1", 2, "value3")

	assert.True(result == expected)
}

func TestVariadic3(t *testing.T) {
	assert := require.New(t)

	si := SampleInterfaceImpl{}

	expected := "value1, value2, value3"
	result := si.Variadic3("value1", "value2", "value3")

	assert.True(result == expected)
}

func TestVariadic4(t *testing.T) {
	assert := require.New(t)

	si := SampleInterfaceImpl{}

	expected := "\"value1\", 2, \"value3\""
	result := si.Variadic4("value1", 2, "value3")

	assert.True(result == expected)
}

func TestMockUnnamed(t *testing.T) {
	assert := require.New(t)

	m := MockSampleInterface{}
	m.On("Unnamed",
		"string", 1, mock.Anything,
	).Return(nil)

	assert.True(m.Unnamed("string", 1, nil) == nil)
}

func TestMockVariadic(t *testing.T) {
	assert := require.New(t)

	m := MockSampleInterface{}
	m.On("Variadic",
		"string1: %s, string2: %s, string3: %s",
		"value1", "value2", "value3",
	).Return("success")

	assert.True(m.Variadic(
		"string1: %s, string2: %s, string3: %s",
		"value1", "value2", "value3",
	) == "success")
}

func TestMockVariadic2(t *testing.T) {
	assert := require.New(t)

	m := MockSampleInterface{}
	m.On("Variadic2",
		"string1: %s, integer2: %d, string3: %s",
		"value1", 2, "value3",
	).Return("success")

	assert.True(m.Variadic2(
		"string1: %s, integer2: %d, string3: %s",
		"value1", 2, "value3",
	) == "success")
}

func TestMockVariadic3(t *testing.T) {
	assert := require.New(t)

	m := MockSampleInterface{}
	m.On("Variadic3",
		"value1", "value2", "value3",
	).Return("success")

	assert.True(m.Variadic3(
		"value1", "value2", "value3",
	) == "success")
}

func TestMockVariadic4(t *testing.T) {
	assert := require.New(t)

	m := MockSampleInterface{}
	m.On("Variadic4",
		"value1", 2, "value3",
	).Return("success")

	assert.True(m.Variadic4(
		"value1", 2, "value3",
	) == "success")
}

func TestMockCallFooBar(t *testing.T) {
	assert := require.New(t)

	m := MockSampleInterface{}
	m.On("CallFooBar",
		mock.Anything,
	).Return(
		func(f foo.Foo) bool {
			return f.Bar()
		},
		"mocked")

	f := mockFoo{}
	f.On("Bar").Return(true)

	b, s := m.CallFooBar(&f)
	assert.True(b == true)
	assert.True(s == "mocked")
	f.AssertNumberOfCalls(t, "Bar", 1)
}

func TestMockCollapsedParams(t *testing.T) {
	assert := require.New(t)

	m := MockSampleInterface{}
	m.On("CollapsedParams",
		mock.Anything,
		mock.Anything,
	).Return("mocked")

	s := m.CollapsedParams([]byte("param1"), []byte("param2"))
	assert.True(s == "mocked")
}

func TestMockCollapsedReturns(t *testing.T) {
	assert := require.New(t)

	m := MockSampleInterface{}
	m.On("CollapsedReturns").Return(
		func() int {
			return 100
		},
		200,
		"mocked",
	)

	x, y, s := m.CollapsedReturns()
	assert.True(x == 100)
	assert.True(y == 200)
	assert.True(s == "mocked")
}

func TestMockVoidReturn(t *testing.T) {
	assert := require.New(t)

	m := MockSampleInterface{}
	m.On("VoidReturn").Once()

	m.VoidReturn()
	assert.True(m.AssertNumberOfCalls(t, "VoidReturn", 1))
}

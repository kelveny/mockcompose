# mockcompose

It seems to be quite difficult in Go to do something which is relatively easy in other object oriented languages. For example, in Java, we can mix real method call with mocked sibling method calls like this:
```java
FooService fooService = PowerMockito.mock(FooService.class);
PowerMockito.doCallRealMethod().when(fooService).SomeFooMethod());
```
In the example, SomeFooMethod() will be called to run real implementation code, while any sibling method that SomeFooMethod() calls will be taken from the mocked version. This ability can give us a fine-grained control in unit testing.

Without such granularity control, we sometimes have to rely heavily on end to end integration tests for Go class implementation. The integration test approach involves a lot of effort to bring up either real or simulated runtime with third-party dependencies, it is also hard and tedious to perform fault injections.

## Solution
[mockery](https://github.com/vektra/mockery) is a very nice tool that allows developers to generate mocking interface implementation in Go, however, it is not enough to satisfy the requriements mentioned above. `mockcompose` can be used to fill the gap, it generates mocking implementation for Go classes, interfaces and functions.

## Install
```
go install github.com/kelveny/mockcompose
```

## Usage

```
  mockcompose can be launched with following options

  -c string
        name of the source class to generate against
  -i string
        name of the source interface to generate against
  -mock value
        name of the function to be mocked
  -n string
        name of the generated class
  -p string
        path of the source package in which to search interfaces and functions
  -pkg string
        name of the package that the generated class resides
  -real value
        name of the method function to be cloned from source class
  -testonly
        if set, append _test to generated file name (default true)
  -v    if set, print verbose logging messages
  -version
        if set, print version information
```
`-pkg` option is usually omitted, `mockcompose` will derive Go package name automatically from current working directory.

You can use multiple `-real` and `-mock` options to specify a set of real class method functions to clone and another set of class method functions to mock.

`mockcompose` is recommended to be used in `go generate`:
```go
//go:generate mockcompose -n testFoo -c foo -real Foo -mock Bar
```
In the example, `mockcompose` will generate a testFoo class with Foo() method function be cloned from real foo class implementation, and Bar() method function be mocked.

source Go class code: foo.go
```go
package foo

type Foo interface {
	Foo() string
	Bar() bool
}

type foo struct {
	name string
}

var _ Foo = (*foo)(nil)

func (f *foo) Foo() string {
	if f.Bar() {
		return "Overriden with Bar"
	}

	return f.name
}

func (f *foo) Bar() bool {
	if f.name == "bar" {
		return true
	}

	return false
}
```

`go generate` configuration: mocks.go
```go
//go:generate mockcompose -n testFoo -c foo -real Foo -mock Bar
//go:generate mockcompose -n FooMock -i Foo
package foo
```

`mockcompose` generated code: mockc_testFoo_test.go
```go
//
// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
//
package foo

import (
	"github.com/stretchr/testify/mock"
)

type testFoo struct {
	foo
	mock.Mock
}

func (f *testFoo) Foo() string {
	if f.Bar() {
		return "Overriden with Bar"
	}
	return f.name
}

func (m *testFoo) Bar() bool {

	_mc_ret := m.Called()

	var _r0 bool

	if _rfn, ok := _mc_ret.Get(0).(func() bool); ok {
		_r0 = _rfn()
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(bool)
		}
	}

	return _r0

}
```

You can now write unit tests to test at fine-grained granularity. This can enable to test individual or a group of class method functions, with dependency closure be mocked.

```go
func TestFoo(t *testing.T) {
    assert := require.New(t)

    fooObj := &testFoo{}

    // Mock sibling method Bar()
    fooObj.On("Bar").Return(false)

    s := fooObj.Foo()
    assert.True(s == "")
}
```

When a class method has external callouts to imported functions from other packages, `mockcompose` also offers function level mock generation. Please check out [mockfn](https://github.com/kelveny/mockcompose/tree/main/test/mockfn) for details.

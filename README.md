# mockcompose

`mockcompose` was originally built to address a Go anti-pattern use case scenario. To be exact, the use case can be described with following Java example:

in Java, we can mix real method call with mocked sibling method calls like this:

```java
FooService fooService = PowerMockito.mock(FooService.class);
PowerMockito.doCallRealMethod().when(fooService).SomeFooMethod());
```

In the example, SomeFooMethod() will be called to run real implementation code, while any sibling method that SomeFooMethod() calls will be taken from the mocked version. This ability can give us fine-grained control in unit testing, in world of Object Oriented languages.

Go is a first-class function programming language, Go best practices prefer small interfaces, in the extreme side of the spectrum, per-function interface would eliminate the needs of such usage pattern to be supported at all in mocking. This might be the reason why most Go mocking tools support only interface mocking.

Nevertheless, if you ever come to here, you may be struggling in balancing the ideal world and practical world, try `mockcompose` to solve your immediate needs and you are recommended to follow Go best practices to refactor your code later, to avoid Go anti-pattern as mentioned above if possible.

`mockcompose` also supports generating [mockery](https://github.com/vektra/mockery) compatible code for Go interfaces and regular functions, which could help pave the way for your code to evolve into ideal shape.

Note: Go class here refers to `Go struct` with functions that take receiver objects of the `struct` type.

## Install

```bash
go install github.com/kelveny/mockcompose
```

## Usage

```text
mockcompose generates mocking implementation for Go classes, interfaces and functions.
  -c string
        name of the source class to generate against
  -help
        if set, print usage information
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
        name of the method function to be cloned from source class or source function
  -testonly
        if set, append _test to generated file name (default true)
  -v    if set, print verbose logging messages
  -version
        if set, print version information
```

`-pkg` option is usually omitted, `mockcompose` will derive Go package name automatically from current working directory.

You can use multiple `-real` and `-mock` options to specify a set of real class method functions to clone and another set of class method functions to mock. The cloned and mocked method functions usually form a test closure. However, in most of cases, it is more convenient to do it at per-method basis. In this way, we clone a method function to test against, have all its callee functions be mocked. The mocked callee closure can be specified in format of `[(this|.|<pkg>)][:(this|.|<pkg>)]*`.

- `this` means to mock all peer callee methods
- `.` means to mock all callee functions that are within the same package as of the testing method
- `<pkg>` means to mock all callee functions from the `<pkg>` package

`mockcompose` is recommended to be used in `go generate`:

```go
//go:generate mockcompose -n testFoo -c foo -real Foo,this:.:fmt
```

In the example, `mockcompose` will generate a `testFoo` class with `Foo()` method function be cloned from real foo class implementation, all callee functions (from package `.` and package `fmt`) and callee peer methods ( indicated by `this`) will be mocked.

source Go class code: `foo.go`

```go
package foo

import "fmt"

type Foo interface {
    Foo() string
    Bar() bool
}

type foo struct {
    name string
}

var _ Foo = (*foo)(nil)
var _ Foo = (*dummyFoo)(nil)

func (f *foo) Foo() string {
    if f.Bar() {
        return "Overriden with Bar"
    }

    dummy()
    fmt.Print("Foo")

    return f.name
}

func (f *foo) Bar() bool {
    return f.name == "bar"
}

```

`go generate` configuration: mocks.go

```go
//go:generate mockcompose -n testFoo -c foo -real Foo,this:.:fmt
//go:generate mockcompose -n FooMock -i Foo
package foo
```

`mockcompose` generated code: mockc_testFoo_test.go

```go
type testFoo struct {
    foo
    mock.Mock
    mock_testFoo_Foo_foo    // named after mock_<test class name>_<test nethod name>_<callee package name>
    mock_testFoo_Foo_fmt
}

type mock_testFoo_Foo_foo struct {
    mock.Mock
}

type mock_testFoo_Foo_fmt struct {
    mock.Mock
}

func (f *testFoo) Foo() string {
    dummy := f.mock_testFoo_Foo_foo.dummy
    fmt := &f.mock_testFoo_Foo_fmt

    if f.Bar() {
        return "Overriden with Bar"
    }
    dummy()
    fmt.Print("Foo")
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

func (m *mock_testFoo_Foo_foo) dummy() {

    m.Called()

}

func (m *mock_testFoo_Foo_fmt) Print(a ...interface{}) (n int, err error) {

    _mc_ret := m.Called(a...)

    var _r0 int

    if _rfn, ok := _mc_ret.Get(0).(func(...interface{}) int); ok {
        _r0 = _rfn(a...)
    } else {
        if _mc_ret.Get(0) != nil {
            _r0 = _mc_ret.Get(0).(int)
        }
    }

    var _r1 error

    if _rfn, ok := _mc_ret.Get(1).(func(...interface{}) error); ok {
        _r1 = _rfn(a...)
    } else {
        _r1 = _mc_ret.Error(1)
    }

    return _r0, _r1

}

```

You can now write unit tests to test at fine-grained granularity.

```go
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
```

## FAQ

### 1. Can mockcompose generate mocked implementation for interfaces?

### __Answer__: Check out `mockcompose` self-test example [mockintf](https://github.com/kelveny/mockcompose/tree/main/test/mockintf)


`go generate` configuration: mocks.go

```go
//go:generate mockcompose -n MockSampleInterface -i SampleInterface
//go:generate mockcompose -n mockFoo -i Foo -p github.com/kelveny/mockcompose/test/foo
package mockintf
```

With this configuration, `mockcompose` generates mocked interface implementation both for an interface defined in its own package and an interface defined in other package.

intf_test.go

```go
package mockintf

import (
    "testing"

    "github.com/kelveny/mockcompose/test/foo"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

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

...

```

### 2. How do I configure `go generate` in YAML?

### __Answer__: Check out `mockcompose` self-test example [yaml](https://github.com/kelveny/mockcompose/tree/main/test/yaml)


`go generate` configuration: mocks.go

```go
//go:generate mockcompose
package yaml
```

`go generate` YAML configuration file: .mockcompose.yaml

```yaml
mockcompose:
  - name: testFoo
    testOnly: true
    className: foo
    real:
      - "Foo,this:.:fmt"
  - name: MockSampleInterface
    testOnly: true
    interfaceName: SampleInterface
  - name: mockFoo
    testOnly: true
    interfaceName: Foo
    sourcePkg: github.com/kelveny/mockcompose/test/foo
```

If `mockcompose` detects `.mockcompose.yaml` or `.mockcompose.yml` in package directory, it will load code generation configuration from the file.

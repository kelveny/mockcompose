# mockcompose

`mockcompose` was originally built to address a Go anti-pattern use case scenario. To be exact, the use case can be described with following Java example:

in Java, we can mix real method call with mocked sibling method calls like this:

```java
FooService fooService = PowerMockito.mock(FooService.class);
PowerMockito.doCallRealMethod().when(fooService).SomeFooMethod());
```

In the example, `SomeFooMethod()` will run its real implementation code, while any sibling methods that `SomeFooMethod()` calls will be taken from the mocked version. This capability provides fine-grained control in unit testing within the realm of object-oriented languages.

Go is a first-class function programming language, and Go best practices favor small interfaces. At the extreme end of the spectrum, using a per-function interface eliminates the need for such usage patterns in mocking. This might be why most Go mocking tools support only interface mocking.

Nevertheless, if you find yourself here, you may be struggling to balance ideal practices with practical needs. Try `mockcompose` to address your immediate requirements, but it's recommended to follow Go best practices and refactor your code later to avoid the aforementioned Go anti-patterns whenever possible.

`mockcompose` also supports generating [mockery](https://github.com/vektra/mockery) compatible code for Go interfaces and regular functions, which can help guide your code toward an ideal structure.

Note: In this context, a Go class refers to a Go `struct` with methods that have receiver objects of the `struct` type.

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

You can use multiple `-real` and `-mock` options to specify a set of functions to clone and aother set to mock. The cloned and mocked functions usually form a __test closure__. However, in most cases, it is more convenient to handle this on a per-method basis. This way, you can clone a function to test against while mocking all its callee functions. The format for specifying the __callee closure__ for automatic code generation is `[(this|.|<pkg>)][:(this|.|<pkg>)]*`.

- `this` means to mock all peer callee methods. Only callees with the same receiver type (either `by-value` or `by-reference`) will be considered as peer callee method.
- `.` means to mock all callee functions that are within the same package as of the function.
- `<pkg>` means to mock all callee functions from the `<pkg>` package. Note, __when you have both references to functions and types from the `<pkg>` package, reference to these functions and types through different import names__.

All mocked function are generated with a `pointer` receiver type. It is also recommended to use `mockcompose` for class with methods that have `pointer` receiver types.

Although mockcompose supports `YAML`-based configuration, in most cases, you may find it more convenient to use `mockcompose` inline with the `//go:generate mockcompose` directive.

## Use cases

### 1. Use `mockcompose` on `per-method` basis

`mockcompose` uses this pattern by itself and follows the convention:

- name the gnerated class in format of `<a shortened version of the class name from the source class>_<method name>` with `-n` option
- use different import names for functions and types from `gosyntax` package

source content (`cmd/clzgenerator.go`):

```go
package cmd

import (
    ...

    "github.com/kelveny/mockcompose/pkg/gosyntax"
    gosyntaxtyp "github.com/kelveny/mockcompose/pkg/gosyntax"

    ...
)

//go:generate mockcompose -n gctx_findClassMethods -c generatorContext -real findClassMethods,gosyntax
func (c *generatorContext) findClassMethods(
    clzTypeDeclString string,
    fset *token.FileSet,
    f *ast.File,
) map[string]*gosyntaxtyp.ReceiverSpec {
    if c.clzMethods == nil {
        c.clzMethods = make(map[string]map[string]*gosyntaxtyp.ReceiverSpec)
    }

    if _, ok := c.clzMethods[clzTypeDeclString]; !ok {
        c.clzMethods[clzTypeDeclString] = gosyntax.FindClassMethods(clzTypeDeclString, fset, f)
    }

    return c.clzMethods[clzTypeDeclString]
}
```

`mockcompose` generated content (`cmd/mockc_gctx_findClassMethods_test.go`):

```go
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
```

test content (`clzgenerator_test.go`):

```go
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

```

### 2. Use `mockcompose` to form a test closure

`mockcompose` directive to generate the closure:

```go
//go:generate mockcompose -n fooBarMock -c fooBar -real FooBar,this -real BarFoo,this:.
```

source content (`bar.go`):

```go
package bar

import (
    "fmt"
    "math/rand"
    "time"
)

type fooBar struct {
    name string
}

//go:generate mockcompose -n fooBarMock -c fooBar -real FooBar,this -real BarFoo,this:.

func (f *fooBar) FooBar() string {
    if f.order()%2 == 0 {
        fmt.Printf("ordinal order\n")

        return fmt.Sprintf("%s: %s%s", f.name, f.Foo(), f.Bar())
    }

    fmt.Printf("reverse order\n")
    return fmt.Sprintf("%s: %s%s", f.name, f.Bar(), f.Foo())
}

func (f *fooBar) BarFoo() string {
    if order()%2 == 0 {
        fmt.Printf("ordinal order\n")

        return fmt.Sprintf("%s: %s%s", f.name, f.Bar(), f.Foo())
    }

    fmt.Printf("reverse order\n")
    return fmt.Sprintf("%s: %s%s", f.name, f.Foo(), f.Bar())
}

func (f *fooBar) Foo() string {
    return "Foo"
}

func (f *fooBar) Bar() string {
    return "Bar"
}

func (f *fooBar) order() int {
    rand.Seed(time.Now().UnixNano())
    return rand.Int()
}

func order() int {
    rand.Seed(time.Now().UnixNano())
    return rand.Int()
}
```

test content:

```go
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
```

### 3. Use `mockcompose` to generate the mocking implementation of a Go interface

`mockcompose` directive to generate for interface `Foo` defined in the same package:

```go
//go:generate mockcompose -n FooMock -i Foo
package foo
```

If the Go interface is defined in external package, specify the `import` path of the package as example:

```go
//go:generate mockcompose -n FooMock -i Foo -sourcePkg github.com/kelveny/mockcompose/test/foo
```

generated implementation of interface `Foo`:

```go
// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package foo

import (
    "github.com/stretchr/testify/mock"
)

type FooMock struct {
    mock.Mock
}

func (m *FooMock) Foo() string {

    _mc_ret := m.Called()

    var _r0 string

    if _rfn, ok := _mc_ret.Get(0).(func() string); ok {
        _r0 = _rfn()
    } else {
        if _mc_ret.Get(0) != nil {
            _r0 = _mc_ret.Get(0).(string)
        }
    }

    return _r0

}

func (m *FooMock) Bar() bool {

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

### 4. Use `mockcompose` for ordinary function

source content:

```go
//go:generate mockcompose -n mockCallee -real functionThatUsesFunctionFromSameRoot,foo
func functionThatUsesFunctionFromSameRoot() string {

    if useRemoteDummy() {
        s := foo.Dummy()
        fmt.Printf("result from remote: %s\n", s)
        return s
    } else {
        s := dummy()

        fmt.Printf("result from local: %s\n", s)
        return s
    }
}

func useRemoteDummy() bool {
    return true
}

func dummy() string {
    return "dummy"
}
```

`mockcompose` will then generate code as:

```go
// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package clonefn

import (
    "fmt"

    "github.com/stretchr/testify/mock"
)

type mockCallee struct {
    mock.Mock
    mock_mockCallee_functionThatUsesFunctionFromSameRoot_foo
}

type mock_mockCallee_functionThatUsesFunctionFromSameRoot_foo struct {
    mock.Mock
}

func (m *mockCallee) functionThatUsesFunctionFromSameRoot() string {
    foo := &m.mock_mockCallee_functionThatUsesFunctionFromSameRoot_foo

    if useRemoteDummy() {
        s := foo.Dummy()
        fmt.Printf("result from remote: %s\n", s)
        return s
    } else {
        s := dummy()
        fmt.Printf("result from local: %s\n", s)
        return s
    }
}

func (m *mock_mockCallee_functionThatUsesFunctionFromSameRoot_foo) Dummy() string {

    _mc_ret := m.Called()

    var _r0 string

    if _rfn, ok := _mc_ret.Get(0).(func() string); ok {
        _r0 = _rfn()
    } else {
        if _mc_ret.Get(0) != nil {
            _r0 = _mc_ret.Get(0).(string)
        }
    }

    return _r0

}

```

test `functionThatUsesFunctionFromSameRoot` with mocked callees:

```go
func Test_functionThatUsesFunctionFromSameRoot(t *testing.T) {
    assert := require.New(t)

    c := &mockCallee{}
    c.mock_mockCallee_functionThatUsesFunctionFromSameRoot_foo.On("Dummy").Return("Mocked Dummy")
    s := c.functionThatUsesFunctionFromSameRoot()
    assert.Equal("Mocked Dummy", s)
}

```

### 5. Configure with `YAML` configuration

If `mockcompose` detects a `.mockcompose.yaml` or `.mockcompose.yml` file in the package directory, it will load the code generation configuration from that file.

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

## Best pratices

- use `mockcompose` for class with methods that have `pointer` receiver types
- use different import names for functions and types from an external package
- for `per-method` basis usage, name the gnerated class in format of `<a shortened version of the class name from the source class>_<method name>`
- for `test-closure` usage, name the generated class in format of `<a shortened version of the class name from the source class>_<a testing aspect derived closure name>`
- be cautious when mocking functions that accept parameters with fields requiring protection in multi-threaded contexts. Inside the [testify implementation](https://github.com/stretchr/testify/blob/master/mock/mock.go#L950), it reads the passed-in parameters without any synchronization on those parameters. Since it doesn't have awareness of the internal concurrency requirements of the object, this can lead to data race conditions, which may be detected by running `go test -race`
  - data-race example: https://github.com/kelveny/mockcompose/blob/main/test/race/race_test.go
  - best practice for data-race avoidance: https://github.com/kelveny/mockcompose/blob/main/test/race2/race_test.go

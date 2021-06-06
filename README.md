# mockcompose

It seems to be quite difficult in Go to do something which is relatively easy in other object oriented languages. For example, in Java, we can mix real method call with mocked sibling method calls like this:
```java
FooService fooService = PowerMockito.mock(FooService.class);
PowerMockito.doCallRealMethod().when(fooService).SomeFooMethod());
```
In the example, SomeFooMethod() will be called to run real implementation code, while any sibling method that SomeFooMethod() calls will be taken from the mocked version. This ability can give us fine-grained control in unit testing.

Without such granularity control, we sometimes have to rely heavily on end to end integration tests for Go class implementation. The integration test approach involves a lot of effort to bring up either real or simulated runtime with third-party dependencies, it is also hard and tedious to perform fault injections.

## Solution
[mockery](https://github.com/vektra/mockery) is a very nice tool that allows developers to generate mocking interface implementation in Go, however, it is not enough to satisfy the requriements mentioned above. `mockcompose` can be used to fill the gap, it generates mocking implementation for Go classes, interfaces and functions.

## Install
```
go install github.com/kelveny/mockcompose
```

## Usage

```
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

## FAQ

### 1. My class method not only has callouts to sibling methods, but also callouts to functions imported from other packages, and I want to mock these imported functions, how can I do that?  
<br/>

### __Answer__: Check out `mockcompose` self-test example [mockfn](https://github.com/kelveny/mockcompose/tree/main/test/mockfn)
<br/>

`go generate` configuration: mocks.go
```go
//go:generate mockcompose -n mockFmt -p fmt -mock Sprintf
//go:generate mockcompose -n mockJson -p encoding/json -mock Marshal
//go:generate mockcompose -n mockSampleClz -c sampleClz -real "methodThatUsesGlobalFunction,fmt=fmtMock"
//go:generate mockcompose -n mockSampleClz2 -c sampleClz -real "methodThatUsesMultileGlobalFunctions,fmt=fmtMock:json=jsonMock"
//go:generate mockcompose -n mockSampleClz3 -c sampleClz -real "methodThatUsesMultileGlobalFunctions,fmt=fmtMock"
package mockfn
```
With this configuration, `mockcompose` generates Go classes for package `fmt` and `encoding/json`, the generated Go classes are equipped with mocked function implementation. `mockcompose` also clones the subject class method with local overrides, thus enables callouts to be redirected to mocked implementation.

fn_test.go
```go
package mockfn

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var jsonMock *mockJson = &mockJson{}
var fmtMock *mockFmt = &mockFmt{}

func TestSampleClz(t *testing.T) {
	assert := require.New(t)

	// setup function mocks
	jsonMock.On("Marshal", mock.Anything).Return(([]byte)("mocked Marshal"), nil)
	fmtMock.On("Sprintf", mock.Anything, mock.Anything).Return("mocked Sprintf")

	// inside mockSampleClz.methodThatUsesMultileGlobalFunctions: fmt.Sprintf is mocked
	sc := mockSampleClz{}
	assert.True(sc.methodThatUsesGlobalFunction("format", "value") == "mocked Sprintf")

	// inside mockSampleClz2.methodThatUsesMultileGlobalFunctions: both json.Marshal()
	// and fmt.Sprintf are mocked
	sc2 := mockSampleClz2{}
	assert.True(sc2.methodThatUsesMultileGlobalFunctions("format", "value") == "mocked Marshalmocked Sprintf")

	// inside mockSampleClz3.methodThatUsesMultileGlobalFunctions: json.Marshal() is not mocked,
	// fmt.Sprintf is mocked
	sc3 := mockSampleClz3{}
	assert.True(sc3.methodThatUsesMultileGlobalFunctions("format", "value") == "\"format\"mocked Sprintf")
}
```

### 2. I want to test my function with callouts to functions imported from other packages, and I want to mock these imported functions, how can I do that?
<br/>

### __Answer__: Check out `mockcompose` self-test example [clonefn](https://github.com/kelveny/mockcompose/tree/main/test/clonefn)
<br/>

`go generate` configuration: mocks.go
```
//go:generate mockcompose -n mockFmt -p fmt -mock Sprintf
//go:generate mockcompose -n mockJson -p encoding/json -mock Marshal
//go:generate mockcompose -n clonedFuncs -real "functionThatUsesMultileGlobalFunctions,fmt=fmtMock:json=jsonMock" -real "functionThatUsesGlobalFunction,fmt=fmtMock" -real "functionThatUsesMultileGlobalFunctions2,fmt=fmtMock"
package clonefn
```
With this configuration, `mockcompose` generates Go classes for package `fmt` and `encoding/json`, the generated Go classes are equipped with mocked function implementation. `mockcompose` also clones the subject function with local overrides, thus enables callouts to be redirected to mocked implementation.

fn_test.go
```go
package clonefn

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var jsonMock *mockJson = &mockJson{}
var fmtMock *mockFmt = &mockFmt{}

func TestClonedFuncs(t *testing.T) {
	assert := require.New(t)

	// setup function mocks
	jsonMock.On("Marshal", mock.Anything).Return(([]byte)("mocked Marshal"), nil)
	fmtMock.On("Sprintf", mock.Anything, mock.Anything).Return("mocked Sprintf")

	// inside functionThatUsesMultileGlobalFunctions: fmt.Sprintf is mocked
	assert.True(functionThatUsesGlobalFunction_clone("format", "value") == "mocked Sprintf")

	// inside functionThatUsesMultileGlobalFunctions: both json.Marshal()
	// and fmt.Sprintf are mocked
	assert.True(functionThatUsesMultileGlobalFunctions_clone("format", "value") == "mocked Marshalmocked Sprintf")

	// inside functionThatUsesMultileGlobalFunctions2: json.Marshal() is not mocked,
	// fmt.Sprintf is mocked
	assert.True(functionThatUsesMultileGlobalFunctions2_clone("format", "value") == "\"format\"mocked Sprintf")
}
```

### 3. Can mockcompose generate mocked implementation for interfaces?
<br/>

### __Answer__: Check out `mockcompose` self-test example [mockintf](https://github.com/kelveny/mockcompose/tree/main/test/mockintf)
<br/>

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

### 4. Can mockcompose generate mocked implementation for functions?
<br/>

### __Answer__: Yes. `mockcompose` can group a set of functions into a generated Go class, the generated Go class has embedded mock object through which function behavior can be mocked.

`go generate` configuration: mocks.go
```go
//go:generate mockcompose -n mockFmt -p fmt -mock Sprintf
//go:generate mockcompose -n mockJson -p encoding/json -mock Marshal
```
With this configuration, `mockcompose` can generate mocking Go class `mockFmt` and `mockJson` that implement `Sprintf` and `Marshal` respectively. Callers of these functions can then use method/function local overrides to connect callouts of method/function to these generated Go classes.

These techniques have been used in examples of the questions above.

<br/>

### 5. How do I configure `go generate` in YAML?
<br/>

### __Answer__: Check out `mockcompose` self-test example [yaml](https://github.com/kelveny/mockcompose/tree/main/test/yaml)
<br/>

`go generate` configuration: mocks.go
```go
//go:generate mockcompose
package yaml
```

`go generate` YAML configuration file: .mockcompose.yaml
```yaml
mockcompose:
  - name: mockFmt
    testOnly: true
    sourcePkg: fmt
    mock: 
      - Sprintf
  - name: mockJson
    testOnly: true
    sourcePkg: encoding/json
    mock: 
      - Marshal
  - name: mockSampleClz
    testOnly: true
    className: sampleClz
    real:
      - "methodThatUsesGlobalFunction,fmt=fmtMock"
  - name: mockSampleClz2
    testOnly: true
    className: sampleClz
    real:
      - "methodThatUsesMultileGlobalFunctions,fmt=fmtMock:json=jsonMock"
  - name: mockSampleClz3
    testOnly: true
    className: sampleClz
    real:
      - "methodThatUsesMultileGlobalFunctions,fmt=fmtMock"
  - name: MockSampleInterface
    testOnly: true
    interfaceName: SampleInterface
  - name: mockFoo
    testOnly: true
    interfaceName: Foo
    sourcePkg: github.com/kelveny/mockcompose/test/foo
  - name: mockFmtclonedFuncs
    testOnly: true
    real: 
      - "functionThatUsesMultileGlobalFunctions,fmt=fmtMock:json=jsonMock" 
      - "functionThatUsesGlobalFunction,fmt=fmtMock" 
      - "functionThatUsesMultileGlobalFunctions2,fmt=fmtMock"
```
If `mockcompose` detects `.mockcompose.yaml` or `.mockcompose.yml` in package directory, it will load code generation configuration from the file.


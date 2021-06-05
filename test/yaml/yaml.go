package yaml

import (
	"encoding/json"
	"fmt"

	"github.com/kelveny/mockcompose/test/foo"
)

type SampleInterface interface {
	// test unnamed parameters
	Unnamed(string, int, chan<- string) error

	// test method/function with variadic parameters
	Variadic(format string, args ...string) string
	Variadic2(format string, args ...interface{}) string
	Variadic3(args ...string) string
	Variadic4(args ...interface{}) string

	// test cross-package imports
	CallFooBar(f foo.Foo) (bool, string)

	// test collapsed paramemters
	CollapsedParams(arg1, arg2 []byte) string

	// test collapsed returns
	CollapsedReturns() (x, y int, z string)

	// test void return
	VoidReturn()
}

type sampleClz struct {
}

func (c *sampleClz) methodThatUsesGlobalFunction(
	format string,
	args ...interface{},
) string {
	//
	// skip fansy logic...
	//

	// call out to a global function in fmt package
	return fmt.Sprintf(format, args...)
}

func (c *sampleClz) methodThatUsesMultileGlobalFunctions(
	format string,
	args ...interface{},
) string {
	//
	// skip fansy logic...
	//

	// call out to a global function in fmt package and json package
	b, _ := json.Marshal(format)
	return string(b) + fmt.Sprintf(format, args...)
}

func functionThatUsesGlobalFunction(
	format string,
	args ...interface{},
) string {
	//
	// skip fansy logic...
	//

	// call out to a global function in fmt package
	return fmt.Sprintf(format, args...)
}

func functionThatUsesMultileGlobalFunctions(
	format string,
	args ...interface{},
) string {
	//
	// skip fansy logic...
	//

	// call out to a global function in fmt package and json package
	b, _ := json.Marshal(format)
	return string(b) + fmt.Sprintf(format, args...)
}

func functionThatUsesMultileGlobalFunctions2(
	format string,
	args ...interface{},
) string {
	//
	// skip fansy logic...
	//

	// call out to a global function in fmt package and json package
	b, _ := json.Marshal(format)
	return string(b) + fmt.Sprintf(format, args...)
}

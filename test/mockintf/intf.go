package mockintf

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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

type SampleInterfaceImpl struct{}

var _ SampleInterface = (*SampleInterfaceImpl)(nil)

func (sii *SampleInterfaceImpl) Unnamed(s string, n int, c chan<- string) error {
	if c == nil {
		return errors.New("Invalid arguments")
	}

	fmt.Printf("s: %s, n: %d\n", s, n)

	return nil
}

func (sii *SampleInterfaceImpl) Variadic(format string, args ...string) string {
	castedArgs := make([]interface{}, len(args))
	for i, a := range args {
		castedArgs[i] = a
	}
	return fmt.Sprintf(format, castedArgs...)
}

func (sii *SampleInterfaceImpl) Variadic2(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func (sii *SampleInterfaceImpl) Variadic3(args ...string) string {
	return strings.Join(args, ", ")
}

func (sii *SampleInterfaceImpl) Variadic4(args ...interface{}) string {
	var r []string
	for _, a := range args {
		r = append(r, toJson(a))
	}

	return strings.Join(r, ", ")
}

func (sii *SampleInterfaceImpl) CallFooBar(f foo.Foo) (bool, string) {
	return f.Bar(), "ok"
}

func (sii *SampleInterfaceImpl) CollapsedParams(arg1, arg2 []byte) string {
	return string(arg1) + string(arg2)
}

func (sii *SampleInterfaceImpl) CollapsedReturns() (x, y int, z string) {
	return 1, 2, ""
}

func (sii *SampleInterfaceImpl) VoidReturn() {
}

func toJson(o interface{}) string {
	b, err := json.Marshal(o)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

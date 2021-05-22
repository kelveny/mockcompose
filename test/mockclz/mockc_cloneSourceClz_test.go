//
// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
//
package mockclz

import (
	"errors"
	"fmt"
	"github.com/kelveny/mockcompose/test/foo"
	"github.com/stretchr/testify/mock"
	"strings"
)

type cloneSourceClz struct {
	sourceClz
	mock.Mock
}

func (sc *cloneSourceClz) Unnamed(s string, n int, c chan<- string) (string, error) {
	if c == nil {
		return "", errors.New("Invalid arguments")
	}
	return fmt.Sprintf("s: %s, n: %d", s, n), nil
}

func (sc cloneSourceClz) Unnamed2(s string, n int, c chan<- string) (string, error) {
	if c == nil {
		return "", errors.New("Invalid arguments")
	}
	return fmt.Sprintf("s: %s, n: %d", s, n), nil
}

func (sc *cloneSourceClz) Variadic(format string, args ...string) string {
	castedArgs := make([]interface{}, len(args))
	for i, a := range args {
		castedArgs[i] = a
	}
	return fmt.Sprintf(format, castedArgs...)
}

func (sc *cloneSourceClz) Variadic2(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func (sc *cloneSourceClz) Variadic3(args ...string) string {
	return strings.Join(args, ", ")
}

func (sc *cloneSourceClz) Variadic4(args ...interface{}) string {
	var r []string
	for _, a := range args {
		r = append(r, toJson(a))
	}
	return strings.Join(r, ", ")
}

func (sc *cloneSourceClz) CallFooBar(f foo.Foo) (bool, string) {
	return f.Bar(), "ok"
}

func (sc *cloneSourceClz) CollapsedParams(arg1, arg2 []byte) string {
	return string(arg1) + string(arg2)
}

func (sc *cloneSourceClz) CollapsedReturns() (x, y int, z string) {
	return 1, 2, ""
}

func (sc *cloneSourceClz) VoidReturn() {
}

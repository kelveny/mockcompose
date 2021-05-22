package mockclz

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/kelveny/mockcompose/test/foo"
)

type sourceClz struct{}

func (sc *sourceClz) Unnamed(s string, n int, c chan<- string) (string, error) {
	if c == nil {
		return "", errors.New("Invalid arguments")
	}

	return fmt.Sprintf("s: %s, n: %d", s, n), nil
}

// test different receiver type
func (sc sourceClz) Unnamed2(s string, n int, c chan<- string) (string, error) {
	if c == nil {
		return "", errors.New("Invalid arguments")
	}

	return fmt.Sprintf("s: %s, n: %d", s, n), nil
}

func (sc *sourceClz) Variadic(format string, args ...string) string {
	castedArgs := make([]interface{}, len(args))
	for i, a := range args {
		castedArgs[i] = a
	}
	return fmt.Sprintf(format, castedArgs...)
}

func (sc *sourceClz) Variadic2(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func (sc *sourceClz) Variadic3(args ...string) string {
	return strings.Join(args, ", ")
}

func (sc *sourceClz) Variadic4(args ...interface{}) string {
	var r []string
	for _, a := range args {
		r = append(r, toJson(a))
	}

	return strings.Join(r, ", ")
}

func (sc *sourceClz) CallFooBar(f foo.Foo) (bool, string) {
	return f.Bar(), "ok"
}

func (sc *sourceClz) CollapsedParams(arg1, arg2 []byte) string {
	return string(arg1) + string(arg2)
}

func (sc *sourceClz) CollapsedReturns() (x, y int, z string) {
	return 1, 2, ""
}

func (sc *sourceClz) VoidReturn() {
}

func toJson(o interface{}) string {
	b, err := json.Marshal(o)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

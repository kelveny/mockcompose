package clonefn

import (
	"encoding/json"
	"fmt"

	"github.com/kelveny/mockcompose/test/foo"
)

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

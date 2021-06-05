package clonefn

import (
	"encoding/json"
	"fmt"
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

	// call out to a global function in fmt package and filepath package
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

	// call out to a global function in fmt package and filepath package
	b, _ := json.Marshal(format)
	return string(b) + fmt.Sprintf(format, args...)
}

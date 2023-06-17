// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package yaml

import (
	"encoding/json"
)

func functionThatUsesGlobalFunction_clone(format string, args ...interface{}) string {
	fmt := fmtMock

	return fmt.Sprintf(format, args...)
}

func functionThatUsesMultileGlobalFunctions_clone(format string, args ...interface{}) string {
	fmt := fmtMock
	json := jsonMock

	b, _ := json.Marshal(format)
	return string(b) + fmt.Sprintf(format, args...)
}

func functionThatUsesMultileGlobalFunctions2_clone(format string, args ...interface{}) string {
	fmt := fmtMock

	b, _ := json.Marshal(format)
	return string(b) + fmt.Sprintf(format, args...)
}

package mockfn

import (
	"encoding/json"
	"fmt"
)

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

// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package mockfn

import (
	"encoding/json"

	"github.com/stretchr/testify/mock"
)

type mockSampleClz3 struct {
	sampleClz
	mock.Mock
}

func (c *mockSampleClz3) methodThatUsesMultileGlobalFunctions(format string, args ...interface{}) string {
	fmt := fmtMock

	b, _ := json.Marshal(format)
	return string(b) + fmt.Sprintf(format, args...)
}

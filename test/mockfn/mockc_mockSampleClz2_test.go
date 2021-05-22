//
// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
//
package mockfn

import (
	"github.com/stretchr/testify/mock"
)

type mockSampleClz2 struct {
	sampleClz
	mock.Mock
}

func (c *mockSampleClz2) methodThatUsesMultileGlobalFunctions(format string, args ...interface{}) string {
	json := jsonMock
	fmt := fmtMock

	b, _ := json.Marshal(format)
	return string(b) + fmt.Sprintf(format, args...)
}

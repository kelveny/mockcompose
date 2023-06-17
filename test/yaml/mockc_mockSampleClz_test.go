// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package yaml

import (
	"github.com/stretchr/testify/mock"
)

type mockSampleClz struct {
	sampleClz
	mock.Mock
}

func (c *mockSampleClz) methodThatUsesGlobalFunction(format string, args ...interface{}) string {
	fmt := fmtMock

	return fmt.Sprintf(format, args...)
}

// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package yaml

import (
	"github.com/stretchr/testify/mock"
)

type mockSampleClz struct {
	sampleClz
	mock.Mock
	mock_mockSampleClz_methodThatUsesGlobalFunction_fmt
}

type mock_mockSampleClz_methodThatUsesGlobalFunction_fmt struct {
	mock.Mock
}

func (c *mockSampleClz) methodThatUsesGlobalFunction(format string, args ...interface{}) string {
	fmt := &c.mock_mockSampleClz_methodThatUsesGlobalFunction_fmt

	return fmt.Sprintf(format, args...)
}

func (m *mock_mockSampleClz_methodThatUsesGlobalFunction_fmt) Sprintf(format string, a ...interface{}) string {

	_mc_args := make([]interface{}, 0, 1+len(a))

	_mc_args = append(_mc_args, format)

	for _, _va := range a {
		_mc_args = append(_mc_args, _va)
	}

	_mc_ret := m.Called(_mc_args...)

	var _r0 string

	if _rfn, ok := _mc_ret.Get(0).(func(string, ...interface{}) string); ok {
		_r0 = _rfn(format, a...)
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(string)
		}
	}

	return _r0

}

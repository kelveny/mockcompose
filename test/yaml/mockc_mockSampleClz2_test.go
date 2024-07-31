// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package yaml

import (
	"github.com/stretchr/testify/mock"
)

type mockSampleClz2 struct {
	sampleClz
	mock.Mock
	mock_mockSampleClz2_methodThatUsesMultileGlobalFunctions_fmt
	mock_mockSampleClz2_methodThatUsesMultileGlobalFunctions_json
}

type mock_mockSampleClz2_methodThatUsesMultileGlobalFunctions_fmt struct {
	mock.Mock
}

type mock_mockSampleClz2_methodThatUsesMultileGlobalFunctions_json struct {
	mock.Mock
}

func (c *mockSampleClz2) methodThatUsesMultileGlobalFunctions(format string, args ...interface{}) string {
	fmt := &c.mock_mockSampleClz2_methodThatUsesMultileGlobalFunctions_fmt
	json := &c.mock_mockSampleClz2_methodThatUsesMultileGlobalFunctions_json

	b, _ := json.Marshal(format)
	return string(b) + fmt.Sprintf(format, args...)
}

func (m *mock_mockSampleClz2_methodThatUsesMultileGlobalFunctions_fmt) Sprintf(format string, a ...interface{}) string {

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

func (m *mock_mockSampleClz2_methodThatUsesMultileGlobalFunctions_json) Marshal(v interface{}) ([]byte, error) {

	_mc_ret := m.Called(v)

	var _r0 []byte

	if _rfn, ok := _mc_ret.Get(0).(func(interface{}) []byte); ok {
		_r0 = _rfn(v)
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).([]byte)
		}
	}

	var _r1 error

	if _rfn, ok := _mc_ret.Get(1).(func(interface{}) error); ok {
		_r1 = _rfn(v)
	} else {
		_r1 = _mc_ret.Error(1)
	}

	return _r0, _r1

}

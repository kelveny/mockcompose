// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package clonefn

import (
	"encoding/json"

	"github.com/stretchr/testify/mock"
)

type clonedFuncs struct {
	mock.Mock
	mock_clonedFuncs_functionThatUsesGlobalFunction_fmt
	mock_clonedFuncs_functionThatUsesMultileGlobalFunctions_fmt
	mock_clonedFuncs_functionThatUsesMultileGlobalFunctions_json
	mock_clonedFuncs_functionThatUsesMultileGlobalFunctions2_fmt
}

type mock_clonedFuncs_functionThatUsesGlobalFunction_fmt struct {
	mock.Mock
}

type mock_clonedFuncs_functionThatUsesMultileGlobalFunctions_fmt struct {
	mock.Mock
}

type mock_clonedFuncs_functionThatUsesMultileGlobalFunctions_json struct {
	mock.Mock
}

type mock_clonedFuncs_functionThatUsesMultileGlobalFunctions2_fmt struct {
	mock.Mock
}

func (m *clonedFuncs) functionThatUsesGlobalFunction(format string, args ...interface{}) string {
	fmt := &m.mock_clonedFuncs_functionThatUsesGlobalFunction_fmt

	return fmt.Sprintf(format, args...)
}

func (m *mock_clonedFuncs_functionThatUsesGlobalFunction_fmt) Sprintf(format string, a ...interface{}) string {

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

func (m *clonedFuncs) functionThatUsesMultileGlobalFunctions(format string, args ...interface{}) string {
	fmt := &m.mock_clonedFuncs_functionThatUsesMultileGlobalFunctions_fmt
	json := &m.mock_clonedFuncs_functionThatUsesMultileGlobalFunctions_json

	b, _ := json.Marshal(format)
	return string(b) + fmt.Sprintf(format, args...)
}

func (m *mock_clonedFuncs_functionThatUsesMultileGlobalFunctions_fmt) Sprintf(format string, a ...interface{}) string {

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

func (m *mock_clonedFuncs_functionThatUsesMultileGlobalFunctions_json) Marshal(v interface{}) ([]byte, error) {

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

func (m *clonedFuncs) functionThatUsesMultileGlobalFunctions2(format string, args ...interface{}) string {
	fmt := &m.mock_clonedFuncs_functionThatUsesMultileGlobalFunctions2_fmt

	b, _ := json.Marshal(format)
	return string(b) + fmt.Sprintf(format, args...)
}

func (m *mock_clonedFuncs_functionThatUsesMultileGlobalFunctions2_fmt) Sprintf(format string, a ...interface{}) string {

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

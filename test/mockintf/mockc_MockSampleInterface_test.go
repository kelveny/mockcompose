// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package mockintf

import (
	"github.com/kelveny/mockcompose/test/foo"
	"github.com/stretchr/testify/mock"
)

type MockSampleInterface struct {
	mock.Mock
}

func (m *MockSampleInterface) Unnamed(_a0 string, _a1 int, _a2 chan<- string) error {

	_mc_ret := m.Called(_a0, _a1, _a2)

	var _r0 error

	if _rfn, ok := _mc_ret.Get(0).(func(string, int, chan<- string) error); ok {
		_r0 = _rfn(_a0, _a1, _a2)
	} else {
		_r0 = _mc_ret.Error(0)
	}

	return _r0

}

func (m *MockSampleInterface) Variadic(format string, args ...string) string {

	_mc_args := make([]interface{}, 0, 1+len(args))

	_mc_args = append(_mc_args, format)

	for _, _va := range args {
		_mc_args = append(_mc_args, _va)
	}

	_mc_ret := m.Called(_mc_args...)

	var _r0 string

	if _rfn, ok := _mc_ret.Get(0).(func(string, ...string) string); ok {
		_r0 = _rfn(format, args...)
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(string)
		}
	}

	return _r0

}

func (m *MockSampleInterface) Variadic2(format string, args ...interface{}) string {

	_mc_args := make([]interface{}, 0, 1+len(args))

	_mc_args = append(_mc_args, format)

	for _, _va := range args {
		_mc_args = append(_mc_args, _va)
	}

	_mc_ret := m.Called(_mc_args...)

	var _r0 string

	if _rfn, ok := _mc_ret.Get(0).(func(string, ...interface{}) string); ok {
		_r0 = _rfn(format, args...)
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(string)
		}
	}

	return _r0

}

func (m *MockSampleInterface) Variadic3(args ...string) string {

	_mc_args := make([]interface{}, 0, 0+len(args))

	for _, _va := range args {
		_mc_args = append(_mc_args, _va)
	}

	_mc_ret := m.Called(_mc_args...)

	var _r0 string

	if _rfn, ok := _mc_ret.Get(0).(func(...string) string); ok {
		_r0 = _rfn(args...)
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(string)
		}
	}

	return _r0

}

func (m *MockSampleInterface) Variadic4(args ...interface{}) string {

	_mc_ret := m.Called(args...)

	var _r0 string

	if _rfn, ok := _mc_ret.Get(0).(func(...interface{}) string); ok {
		_r0 = _rfn(args...)
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(string)
		}
	}

	return _r0

}

func (m *MockSampleInterface) CallFooBar(f foo.Foo) (bool, string) {

	_mc_ret := m.Called(f)

	var _r0 bool

	if _rfn, ok := _mc_ret.Get(0).(func(foo.Foo) bool); ok {
		_r0 = _rfn(f)
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(bool)
		}
	}

	var _r1 string

	if _rfn, ok := _mc_ret.Get(1).(func(foo.Foo) string); ok {
		_r1 = _rfn(f)
	} else {
		if _mc_ret.Get(1) != nil {
			_r1 = _mc_ret.Get(1).(string)
		}
	}

	return _r0, _r1

}

func (m *MockSampleInterface) CollapsedParams(arg1 []byte, arg2 []byte) string {

	_mc_ret := m.Called(arg1, arg2)

	var _r0 string

	if _rfn, ok := _mc_ret.Get(0).(func([]byte, []byte) string); ok {
		_r0 = _rfn(arg1, arg2)
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(string)
		}
	}

	return _r0

}

func (m *MockSampleInterface) CollapsedReturns() (x int, y int, z string) {

	_mc_ret := m.Called()

	var _r0 int

	if _rfn, ok := _mc_ret.Get(0).(func() int); ok {
		_r0 = _rfn()
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(int)
		}
	}

	var _r1 int

	if _rfn, ok := _mc_ret.Get(1).(func() int); ok {
		_r1 = _rfn()
	} else {
		if _mc_ret.Get(1) != nil {
			_r1 = _mc_ret.Get(1).(int)
		}
	}

	var _r2 string

	if _rfn, ok := _mc_ret.Get(2).(func() string); ok {
		_r2 = _rfn()
	} else {
		if _mc_ret.Get(2) != nil {
			_r2 = _mc_ret.Get(2).(string)
		}
	}

	return _r0, _r1, _r2

}

func (m *MockSampleInterface) VoidReturn() {

	m.Called()

}

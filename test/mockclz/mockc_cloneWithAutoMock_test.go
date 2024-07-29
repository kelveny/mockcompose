// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package mockclz

import (
	"github.com/stretchr/testify/mock"
)

type cloneWithAutoMock struct {
	sourceClz
	mock.Mock
	mock_cloneWithAutoMock_CallPeer_mockclz
	mock_cloneWithAutoMock_CallPeer_fmt
}

type mock_cloneWithAutoMock_CallPeer_mockclz struct {
	mock.Mock
}

type mock_cloneWithAutoMock_CallPeer_fmt struct {
	mock.Mock
}

func (sc *cloneWithAutoMock) CallPeer() {
	dummy := sc.mock_cloneWithAutoMock_CallPeer_mockclz.dummy
	fmt := &sc.mock_cloneWithAutoMock_CallPeer_fmt
	toJson := sc.mock_cloneWithAutoMock_CallPeer_mockclz.toJson

	sc.Variadic("dummy")
	toJson(fmt.Sprintf("dummy %s", sc.Variadic4()))
	dummy()
	fmt.Printf("dummy")
}

func (m *cloneWithAutoMock) Variadic(format string, args ...string) string {

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

func (m *cloneWithAutoMock) Variadic4(args ...interface{}) string {

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

func (m *mock_cloneWithAutoMock_CallPeer_mockclz) toJson(o interface{}) string {

	_mc_ret := m.Called(o)

	var _r0 string

	if _rfn, ok := _mc_ret.Get(0).(func(interface{}) string); ok {
		_r0 = _rfn(o)
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(string)
		}
	}

	return _r0

}

func (m *mock_cloneWithAutoMock_CallPeer_mockclz) dummy() {

	m.Called()

}

func (m *mock_cloneWithAutoMock_CallPeer_fmt) Sprintf(format string, a ...interface{}) string {

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

func (m *mock_cloneWithAutoMock_CallPeer_fmt) Printf(format string, a ...interface{}) (n int, err error) {

	_mc_args := make([]interface{}, 0, 1+len(a))

	_mc_args = append(_mc_args, format)

	for _, _va := range a {
		_mc_args = append(_mc_args, _va)
	}

	_mc_ret := m.Called(_mc_args...)

	var _r0 int

	if _rfn, ok := _mc_ret.Get(0).(func(string, ...interface{}) int); ok {
		_r0 = _rfn(format, a...)
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(int)
		}
	}

	var _r1 error

	if _rfn, ok := _mc_ret.Get(1).(func(string, ...interface{}) error); ok {
		_r1 = _rfn(format, a...)
	} else {
		_r1 = _mc_ret.Error(1)
	}

	return _r0, _r1

}

// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package mockintf

import (
	"github.com/stretchr/testify/mock"
)

type mockFoo struct {
	mock.Mock
}

func (m *mockFoo) Foo() string {

	_mc_ret := m.Called()

	var _r0 string

	if _rfn, ok := _mc_ret.Get(0).(func() string); ok {
		_r0 = _rfn()
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(string)
		}
	}

	return _r0

}

func (m *mockFoo) Bar() bool {

	_mc_ret := m.Called()

	var _r0 bool

	if _rfn, ok := _mc_ret.Get(0).(func() bool); ok {
		_r0 = _rfn()
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(bool)
		}
	}

	return _r0

}

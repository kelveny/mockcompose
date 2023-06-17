// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package clonefn

import (
	"github.com/stretchr/testify/mock"
)

type mockJson struct {
	mock.Mock
}

func (m *mockJson) Marshal(v interface{}) ([]byte, error) {

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

// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package clonefn

import (
	"fmt"

	"github.com/stretchr/testify/mock"
)

type mockCallee struct {
	mock.Mock
	mock_mockCallee_functionThatUsesFunctionFromSameRoot_foo
}

type mock_mockCallee_functionThatUsesFunctionFromSameRoot_foo struct {
	mock.Mock
}

func (m *mockCallee) functionThatUsesFunctionFromSameRoot() string {
	foo := &m.mock_mockCallee_functionThatUsesFunctionFromSameRoot_foo

	if useRemoteDummy() {
		s := foo.Dummy()
		fmt.Printf("result from remote: %s\n", s)
		return s
	} else {
		s := dummy()
		fmt.Printf("result from local: %s\n", s)
		return s
	}
}

func (m *mock_mockCallee_functionThatUsesFunctionFromSameRoot_foo) Dummy() string {

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

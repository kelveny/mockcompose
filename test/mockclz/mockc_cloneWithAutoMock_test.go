// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package mockclz

import (
	"github.com/stretchr/testify/mock"
)

type cloneWithAutoMock struct {
	sourceClz
	mock.Mock
}

func (sc *cloneWithAutoMock) CallPeer() {
	sc.Variadic("dummy")
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

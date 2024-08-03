// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package mix

import (
	"github.com/stretchr/testify/mock"
)

type mix_checkAndSetOnTarget struct {
	mixReceiver
	mock.Mock
}

func (v mix_checkAndSetOnTarget) checkAndSetOnTarget(p *mixReceiver, s string, val string) {
	if v.getValue() != s {
		v.setValue(val)
		p.setValue(val)
	}
}

func (m *mix_checkAndSetOnTarget) getValue() string {

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

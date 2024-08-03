// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package mix

import (
	"github.com/stretchr/testify/mock"
)

type mix_checkAndSet struct {
	mixReceiver
	mock.Mock
}

func (p *mix_checkAndSet) checkAndSet(s string, val string) {
	if p.getValue() != s {
		p.setValue(val)
	}
}

func (m *mix_checkAndSet) setValue(val string) {

	m.Called(val)

}

// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package foo

import (
	"github.com/stretchr/testify/mock"
)

type testFoo struct {
	foo
	mock.Mock
}

func (f *testFoo) Foo() string {
	if f.Bar() {
		return "Overriden with Bar"
	}
	return f.name
}

func (m *testFoo) Bar() bool {

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

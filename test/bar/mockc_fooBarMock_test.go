// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package bar

import (
	"fmt"

	"github.com/stretchr/testify/mock"
)

type fooBarMock struct {
	fooBar
	mock.Mock
	mock_fooBarMock_BarFoo_bar
}

type mock_fooBarMock_BarFoo_bar struct {
	mock.Mock
}

func (f *fooBarMock) FooBar() string {
	if f.order()%2 == 0 {
		fmt.Printf("ordinal order\n")
		return fmt.Sprintf("%s: %s%s", f.name, f.Foo(), f.Bar())
	}
	fmt.Printf("reverse order\n")
	return fmt.Sprintf("%s: %s%s", f.name, f.Bar(), f.Foo())
}

func (m *fooBarMock) order() int {

	_mc_ret := m.Called()

	var _r0 int

	if _rfn, ok := _mc_ret.Get(0).(func() int); ok {
		_r0 = _rfn()
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(int)
		}
	}

	return _r0

}

func (m *fooBarMock) Foo() string {

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

func (m *fooBarMock) Bar() string {

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

func (f *fooBarMock) BarFoo() string {
	order := f.mock_fooBarMock_BarFoo_bar.order

	if order()%2 == 0 {
		fmt.Printf("ordinal order\n")
		return fmt.Sprintf("%s: %s%s", f.name, f.Bar(), f.Foo())
	}
	fmt.Printf("reverse order\n")
	return fmt.Sprintf("%s: %s%s", f.name, f.Foo(), f.Bar())
}

func (m *mock_fooBarMock_BarFoo_bar) order() int {

	_mc_ret := m.Called()

	var _r0 int

	if _rfn, ok := _mc_ret.Get(0).(func() int); ok {
		_r0 = _rfn()
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(int)
		}
	}

	return _r0

}

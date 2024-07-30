// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
package foo

import (
	"github.com/stretchr/testify/mock"
)

type testFoo struct {
	foo
	mock.Mock
	mock_testFoo_Foo_foo
	mock_testFoo_Foo_fmt
}

type mock_testFoo_Foo_foo struct {
	mock.Mock
}

type mock_testFoo_Foo_fmt struct {
	mock.Mock
}

func (f *testFoo) Foo() string {
	dummy := f.mock_testFoo_Foo_foo.dummy
	fmt := &f.mock_testFoo_Foo_fmt

	if f.Bar() {
		return "Overriden with Bar"
	}
	dummy()
	fmt.Print("Foo")
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

func (m *mock_testFoo_Foo_foo) dummy() {

	m.Called()

}

func (m *mock_testFoo_Foo_fmt) Print(a ...interface{}) (n int, err error) {

	_mc_ret := m.Called(a...)

	var _r0 int

	if _rfn, ok := _mc_ret.Get(0).(func(...interface{}) int); ok {
		_r0 = _rfn(a...)
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(int)
		}
	}

	var _r1 error

	if _rfn, ok := _mc_ret.Get(1).(func(...interface{}) error); ok {
		_r1 = _rfn(a...)
	} else {
		_r1 = _mc_ret.Error(1)
	}

	return _r0, _r1

}

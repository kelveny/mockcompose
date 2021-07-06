//
// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
//
package mockfn

import (
	"regexp"

	"github.com/kelveny/mockcompose/test/libfn"
	"github.com/stretchr/testify/mock"
)

type mockLibfn struct {
	mock.Mock
}

func (m *mockLibfn) GetSecrets(projectId string, secretsRegexp *regexp.Regexp) ([]libfn.SecretData, error) {

	_mc_ret := m.Called(projectId, secretsRegexp)

	var _r0 []libfn.SecretData

	if _rfn, ok := _mc_ret.Get(0).(func(string, *regexp.Regexp) []libfn.SecretData); ok {
		_r0 = _rfn(projectId, secretsRegexp)
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).([]libfn.SecretData)
		}
	}

	var _r1 error

	if _rfn, ok := _mc_ret.Get(1).(func(string, *regexp.Regexp) error); ok {
		_r1 = _rfn(projectId, secretsRegexp)
	} else {
		_r1 = _mc_ret.Error(1)
	}

	return _r0, _r1

}

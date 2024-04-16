// Code generated by mockery v3.0.0-alpha.0. DO NOT EDIT.

package mocks

import (
	jwt "github.com/cristalhq/jwt/v5"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// AuthManager is an autogenerated mock type for the AuthManager type
type AuthManager struct {
	mock.Mock
}

// BuildToken provides a mock function with given fields: userID
func (_m *AuthManager) BuildToken(userID uint) (*jwt.Token, error) {
	ret := _m.Called(userID)

	var r0 *jwt.Token
	if rf, ok := ret.Get(0).(func(uint) *jwt.Token); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jwt.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTokenTTL provides a mock function with given fields:
func (_m *AuthManager) GetTokenTTL() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// HashString provides a mock function with given fields: text
func (_m *AuthManager) HashString(text string) (string, error) {
	ret := _m.Called(text)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(text)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(text)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAuthManager interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthManager creates a new instance of AuthManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthManager(t mockConstructorTestingTNewAuthManager) *AuthManager {
	mock := &AuthManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
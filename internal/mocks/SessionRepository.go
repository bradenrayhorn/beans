// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	beans "github.com/bradenrayhorn/beans/beans"
	mock "github.com/stretchr/testify/mock"
)

// SessionRepository is an autogenerated mock type for the SessionRepository type
type SessionRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: userID
func (_m *SessionRepository) Create(userID beans.UserID) (*beans.Session, error) {
	ret := _m.Called(userID)

	var r0 *beans.Session
	if rf, ok := ret.Get(0).(func(beans.UserID) *beans.Session); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*beans.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(beans.UserID) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: id
func (_m *SessionRepository) Get(id beans.SessionID) (*beans.Session, error) {
	ret := _m.Called(id)

	var r0 *beans.Session
	if rf, ok := ret.Get(0).(func(beans.SessionID) *beans.Session); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*beans.Session)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(beans.SessionID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewSessionRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewSessionRepository creates a new instance of SessionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSessionRepository(t mockConstructorTestingTNewSessionRepository) *SessionRepository {
	mock := &SessionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
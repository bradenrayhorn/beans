// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	beans "github.com/bradenrayhorn/beans/beans"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MonthRepository is an autogenerated mock type for the MonthRepository type
type MonthRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, month
func (_m *MonthRepository) Create(ctx context.Context, month *beans.Month) error {
	ret := _m.Called(ctx, month)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *beans.Month) error); ok {
		r0 = rf(ctx, month)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByDate provides a mock function with given fields: ctx, budgetID, date
func (_m *MonthRepository) GetByDate(ctx context.Context, budgetID beans.ID, date time.Time) (*beans.Month, error) {
	ret := _m.Called(ctx, budgetID, date)

	var r0 *beans.Month
	if rf, ok := ret.Get(0).(func(context.Context, beans.ID, time.Time) *beans.Month); ok {
		r0 = rf(ctx, budgetID, date)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*beans.Month)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, beans.ID, time.Time) error); ok {
		r1 = rf(ctx, budgetID, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMonthRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewMonthRepository creates a new instance of MonthRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMonthRepository(t mockConstructorTestingTNewMonthRepository) *MonthRepository {
	mock := &MonthRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

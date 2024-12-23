// Code generated by mockery v2.46.0. DO NOT EDIT.

package repository

import (
	context "context"

	model "code.emcdtech.com/emcd/service/referral/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// MockDefaultSettings is an autogenerated mock type for the DefaultSettings type
type MockDefaultSettings struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, in
func (_m *MockDefaultSettings) Create(ctx context.Context, in *model.DefaultSettings) error {
	ret := _m.Called(ctx, in)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.DefaultSettings) error); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, product, coin
func (_m *MockDefaultSettings) Delete(ctx context.Context, product string, coin string) error {
	ret := _m.Called(ctx, product, coin)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, product, coin)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, product, coin
func (_m *MockDefaultSettings) Get(ctx context.Context, product string, coin string) (*model.DefaultSettings, error) {
	ret := _m.Called(ctx, product, coin)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *model.DefaultSettings
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*model.DefaultSettings, error)); ok {
		return rf(ctx, product, coin)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.DefaultSettings); ok {
		r0 = rf(ctx, product, coin)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DefaultSettings)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, product, coin)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx, skip, take
func (_m *MockDefaultSettings) GetAll(ctx context.Context, skip int32, take int32) ([]*model.DefaultSettings, int, error) {
	ret := _m.Called(ctx, skip, take)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []*model.DefaultSettings
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, int32, int32) ([]*model.DefaultSettings, int, error)); ok {
		return rf(ctx, skip, take)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int32, int32) []*model.DefaultSettings); ok {
		r0 = rf(ctx, skip, take)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.DefaultSettings)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int32, int32) int); ok {
		r1 = rf(ctx, skip, take)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, int32, int32) error); ok {
		r2 = rf(ctx, skip, take)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetAllWithoutPagination provides a mock function with given fields: ctx
func (_m *MockDefaultSettings) GetAllWithoutPagination(ctx context.Context) ([]*model.DefaultSettings, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllWithoutPagination")
	}

	var r0 []*model.DefaultSettings
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*model.DefaultSettings, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*model.DefaultSettings); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.DefaultSettings)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSettingByReferrer provides a mock function with given fields: ctx, referrerUUID
func (_m *MockDefaultSettings) GetSettingByReferrer(ctx context.Context, referrerUUID string) ([]*model.DefaultSettings, error) {
	ret := _m.Called(ctx, referrerUUID)

	if len(ret) == 0 {
		panic("no return value specified for GetSettingByReferrer")
	}

	var r0 []*model.DefaultSettings
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*model.DefaultSettings, error)); ok {
		return rf(ctx, referrerUUID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.DefaultSettings); ok {
		r0 = rf(ctx, referrerUUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.DefaultSettings)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, referrerUUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, in
func (_m *MockDefaultSettings) Update(ctx context.Context, in *model.DefaultSettings) error {
	ret := _m.Called(ctx, in)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.DefaultSettings) error); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockDefaultSettings creates a new instance of MockDefaultSettings. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDefaultSettings(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDefaultSettings {
	mock := &MockDefaultSettings{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

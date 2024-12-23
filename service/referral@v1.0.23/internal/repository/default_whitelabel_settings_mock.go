// Code generated by mockery v2.46.0. DO NOT EDIT.

package repository

import (
	context "context"

	model "code.emcdtech.com/emcd/service/referral/internal/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockDefaultWhitelabelSettings is an autogenerated mock type for the DefaultWhitelabelSettings type
type MockDefaultWhitelabelSettings struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, in
func (_m *MockDefaultWhitelabelSettings) Create(ctx context.Context, in *model.DefaultWhitelabelSettingsV2) error {
	ret := _m.Called(ctx, in)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.DefaultWhitelabelSettingsV2) error); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, product, coin, whitelabelID
func (_m *MockDefaultWhitelabelSettings) Delete(ctx context.Context, product string, coin string, whitelabelID uuid.UUID) error {
	ret := _m.Called(ctx, product, coin, whitelabelID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, uuid.UUID) error); ok {
		r0 = rf(ctx, product, coin, whitelabelID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, product, coin, whitelabelID
func (_m *MockDefaultWhitelabelSettings) Get(ctx context.Context, product string, coin string, whitelabelID uuid.UUID) (*model.DefaultWhitelabelSettingsV2, error) {
	ret := _m.Called(ctx, product, coin, whitelabelID)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *model.DefaultWhitelabelSettingsV2
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, uuid.UUID) (*model.DefaultWhitelabelSettingsV2, error)); ok {
		return rf(ctx, product, coin, whitelabelID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, uuid.UUID) *model.DefaultWhitelabelSettingsV2); ok {
		r0 = rf(ctx, product, coin, whitelabelID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DefaultWhitelabelSettingsV2)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, uuid.UUID) error); ok {
		r1 = rf(ctx, product, coin, whitelabelID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllWithFilters provides a mock function with given fields: ctx, skip, take, filters
func (_m *MockDefaultWhitelabelSettings) GetAllWithFilters(ctx context.Context, skip int32, take int32, filters map[string]string) ([]*model.DefaultWhitelabelSettingsV2, int, error) {
	ret := _m.Called(ctx, skip, take, filters)

	if len(ret) == 0 {
		panic("no return value specified for GetAllWithFilters")
	}

	var r0 []*model.DefaultWhitelabelSettingsV2
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, int32, int32, map[string]string) ([]*model.DefaultWhitelabelSettingsV2, int, error)); ok {
		return rf(ctx, skip, take, filters)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int32, int32, map[string]string) []*model.DefaultWhitelabelSettingsV2); ok {
		r0 = rf(ctx, skip, take, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.DefaultWhitelabelSettingsV2)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int32, int32, map[string]string) int); ok {
		r1 = rf(ctx, skip, take, filters)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, int32, int32, map[string]string) error); ok {
		r2 = rf(ctx, skip, take, filters)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetAllWithoutPaginationWithFilters provides a mock function with given fields: ctx, filters
func (_m *MockDefaultWhitelabelSettings) GetAllWithoutPaginationWithFilters(ctx context.Context, filters map[string]string) ([]*model.DefaultWhitelabelSettingsV2, error) {
	ret := _m.Called(ctx, filters)

	if len(ret) == 0 {
		panic("no return value specified for GetAllWithoutPaginationWithFilters")
	}

	var r0 []*model.DefaultWhitelabelSettingsV2
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]string) ([]*model.DefaultWhitelabelSettingsV2, error)); ok {
		return rf(ctx, filters)
	}
	if rf, ok := ret.Get(0).(func(context.Context, map[string]string) []*model.DefaultWhitelabelSettingsV2); ok {
		r0 = rf(ctx, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.DefaultWhitelabelSettingsV2)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, map[string]string) error); ok {
		r1 = rf(ctx, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetV2 provides a mock function with given fields: ctx, wlID
func (_m *MockDefaultWhitelabelSettings) GetV2(ctx context.Context, wlID uuid.UUID) ([]*model.DefaultWhitelabelSettingsV2, error) {
	ret := _m.Called(ctx, wlID)

	if len(ret) == 0 {
		panic("no return value specified for GetV2")
	}

	var r0 []*model.DefaultWhitelabelSettingsV2
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]*model.DefaultWhitelabelSettingsV2, error)); ok {
		return rf(ctx, wlID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []*model.DefaultWhitelabelSettingsV2); ok {
		r0 = rf(ctx, wlID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.DefaultWhitelabelSettingsV2)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, wlID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetV2ByCoin provides a mock function with given fields: ctx, product, coin, wlID
func (_m *MockDefaultWhitelabelSettings) GetV2ByCoin(ctx context.Context, product string, coin string, wlID uuid.UUID) (*model.DefaultWhitelabelSettingsV2, error) {
	ret := _m.Called(ctx, product, coin, wlID)

	if len(ret) == 0 {
		panic("no return value specified for GetV2ByCoin")
	}

	var r0 *model.DefaultWhitelabelSettingsV2
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, uuid.UUID) (*model.DefaultWhitelabelSettingsV2, error)); ok {
		return rf(ctx, product, coin, wlID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, uuid.UUID) *model.DefaultWhitelabelSettingsV2); ok {
		r0 = rf(ctx, product, coin, wlID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DefaultWhitelabelSettingsV2)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, uuid.UUID) error); ok {
		r1 = rf(ctx, product, coin, wlID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, in
func (_m *MockDefaultWhitelabelSettings) Update(ctx context.Context, in *model.DefaultWhitelabelSettingsV2) error {
	ret := _m.Called(ctx, in)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.DefaultWhitelabelSettingsV2) error); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockDefaultWhitelabelSettings creates a new instance of MockDefaultWhitelabelSettings. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDefaultWhitelabelSettings(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDefaultWhitelabelSettings {
	mock := &MockDefaultWhitelabelSettings{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

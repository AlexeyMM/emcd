// Code generated by mockery v2.46.3. DO NOT EDIT.

package service

import (
	context "context"

	model "code.emcdtech.com/b2b/swap/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockWithdraw is an autogenerated mock type for the Withdraw type
type MockWithdraw struct {
	mock.Mock
}

type MockWithdraw_Expecter struct {
	mock *mock.Mock
}

func (_m *MockWithdraw) EXPECT() *MockWithdraw_Expecter {
	return &MockWithdraw_Expecter{mock: &_m.Mock}
}

// GetBySwapID provides a mock function with given fields: ctx, swapID
func (_m *MockWithdraw) GetBySwapID(ctx context.Context, swapID uuid.UUID) (*model.Withdraw, error) {
	ret := _m.Called(ctx, swapID)

	if len(ret) == 0 {
		panic("no return value specified for GetBySwapID")
	}

	var r0 *model.Withdraw
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*model.Withdraw, error)); ok {
		return rf(ctx, swapID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *model.Withdraw); ok {
		r0 = rf(ctx, swapID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Withdraw)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, swapID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockWithdraw_GetBySwapID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBySwapID'
type MockWithdraw_GetBySwapID_Call struct {
	*mock.Call
}

// GetBySwapID is a helper method to define mock.On call
//   - ctx context.Context
//   - swapID uuid.UUID
func (_e *MockWithdraw_Expecter) GetBySwapID(ctx interface{}, swapID interface{}) *MockWithdraw_GetBySwapID_Call {
	return &MockWithdraw_GetBySwapID_Call{Call: _e.mock.On("GetBySwapID", ctx, swapID)}
}

func (_c *MockWithdraw_GetBySwapID_Call) Run(run func(ctx context.Context, swapID uuid.UUID)) *MockWithdraw_GetBySwapID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockWithdraw_GetBySwapID_Call) Return(_a0 *model.Withdraw, _a1 error) *MockWithdraw_GetBySwapID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockWithdraw_GetBySwapID_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*model.Withdraw, error)) *MockWithdraw_GetBySwapID_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockWithdraw creates a new instance of MockWithdraw. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockWithdraw(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockWithdraw {
	mock := &MockWithdraw{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
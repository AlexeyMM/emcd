// Code generated by mockery v2.46.3. DO NOT EDIT.

package service

import (
	context "context"

	model "code.emcdtech.com/b2b/swap/model"
	mock "github.com/stretchr/testify/mock"
)

// MockTransfer is an autogenerated mock type for the Transfer type
type MockTransfer struct {
	mock.Mock
}

type MockTransfer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTransfer) EXPECT() *MockTransfer_Expecter {
	return &MockTransfer_Expecter{mock: &_m.Mock}
}

// GetLastInternalTransfer provides a mock function with given fields: ctx, accountID
func (_m *MockTransfer) GetLastInternalTransfer(ctx context.Context, accountID int64) (*model.InternalTransfer, error) {
	ret := _m.Called(ctx, accountID)

	if len(ret) == 0 {
		panic("no return value specified for GetLastInternalTransfer")
	}

	var r0 *model.InternalTransfer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*model.InternalTransfer, error)); ok {
		return rf(ctx, accountID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *model.InternalTransfer); ok {
		r0 = rf(ctx, accountID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.InternalTransfer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, accountID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTransfer_GetLastInternalTransfer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLastInternalTransfer'
type MockTransfer_GetLastInternalTransfer_Call struct {
	*mock.Call
}

// GetLastInternalTransfer is a helper method to define mock.On call
//   - ctx context.Context
//   - accountID int64
func (_e *MockTransfer_Expecter) GetLastInternalTransfer(ctx interface{}, accountID interface{}) *MockTransfer_GetLastInternalTransfer_Call {
	return &MockTransfer_GetLastInternalTransfer_Call{Call: _e.mock.On("GetLastInternalTransfer", ctx, accountID)}
}

func (_c *MockTransfer_GetLastInternalTransfer_Call) Run(run func(ctx context.Context, accountID int64)) *MockTransfer_GetLastInternalTransfer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockTransfer_GetLastInternalTransfer_Call) Return(_a0 *model.InternalTransfer, _a1 error) *MockTransfer_GetLastInternalTransfer_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTransfer_GetLastInternalTransfer_Call) RunAndReturn(run func(context.Context, int64) (*model.InternalTransfer, error)) *MockTransfer_GetLastInternalTransfer_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTransfer creates a new instance of MockTransfer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTransfer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTransfer {
	mock := &MockTransfer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Code generated by mockery v2.46.3. DO NOT EDIT.

package service

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockOrderFee is an autogenerated mock type for the OrderFee type
type MockOrderFee struct {
	mock.Mock
}

type MockOrderFee_Expecter struct {
	mock *mock.Mock
}

func (_m *MockOrderFee) EXPECT() *MockOrderFee_Expecter {
	return &MockOrderFee_Expecter{mock: &_m.Mock}
}

// SyncWithAPI provides a mock function with given fields: ctx
func (_m *MockOrderFee) SyncWithAPI(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for SyncWithAPI")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockOrderFee_SyncWithAPI_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SyncWithAPI'
type MockOrderFee_SyncWithAPI_Call struct {
	*mock.Call
}

// SyncWithAPI is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockOrderFee_Expecter) SyncWithAPI(ctx interface{}) *MockOrderFee_SyncWithAPI_Call {
	return &MockOrderFee_SyncWithAPI_Call{Call: _e.mock.On("SyncWithAPI", ctx)}
}

func (_c *MockOrderFee_SyncWithAPI_Call) Run(run func(ctx context.Context)) *MockOrderFee_SyncWithAPI_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockOrderFee_SyncWithAPI_Call) Return(_a0 error) *MockOrderFee_SyncWithAPI_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockOrderFee_SyncWithAPI_Call) RunAndReturn(run func(context.Context) error) *MockOrderFee_SyncWithAPI_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockOrderFee creates a new instance of MockOrderFee. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockOrderFee(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockOrderFee {
	mock := &MockOrderFee{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

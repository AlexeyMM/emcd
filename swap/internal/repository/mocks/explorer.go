// Code generated by mockery v2.43.2. DO NOT EDIT.

package repository

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockExplorer is an autogenerated mock type for the Explorer type
type MockExplorer struct {
	mock.Mock
}

type MockExplorer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockExplorer) EXPECT() *MockExplorer_Expecter {
	return &MockExplorer_Expecter{mock: &_m.Mock}
}

// GetTransactionLink provides a mock function with given fields: ctx, coin, hashID
func (_m *MockExplorer) GetTransactionLink(ctx context.Context, coin string, hashID string) (string, error) {
	ret := _m.Called(ctx, coin, hashID)

	if len(ret) == 0 {
		panic("no return value specified for GetTransactionLink")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (string, error)); ok {
		return rf(ctx, coin, hashID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, coin, hashID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, coin, hashID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockExplorer_GetTransactionLink_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTransactionLink'
type MockExplorer_GetTransactionLink_Call struct {
	*mock.Call
}

// GetTransactionLink is a helper method to define mock.On call
//   - ctx context.Context
//   - coin string
//   - hashID string
func (_e *MockExplorer_Expecter) GetTransactionLink(ctx interface{}, coin interface{}, hashID interface{}) *MockExplorer_GetTransactionLink_Call {
	return &MockExplorer_GetTransactionLink_Call{Call: _e.mock.On("GetTransactionLink", ctx, coin, hashID)}
}

func (_c *MockExplorer_GetTransactionLink_Call) Run(run func(ctx context.Context, coin string, hashID string)) *MockExplorer_GetTransactionLink_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockExplorer_GetTransactionLink_Call) Return(_a0 string, _a1 error) *MockExplorer_GetTransactionLink_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockExplorer_GetTransactionLink_Call) RunAndReturn(run func(context.Context, string, string) (string, error)) *MockExplorer_GetTransactionLink_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockExplorer creates a new instance of MockExplorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockExplorer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockExplorer {
	mock := &MockExplorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

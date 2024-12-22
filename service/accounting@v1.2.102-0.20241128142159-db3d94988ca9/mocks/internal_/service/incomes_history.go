// Code generated by mockery v2.43.2. DO NOT EDIT.

package service

import (
	context "context"

	accounting "code.emcdtech.com/emcd/service/accounting/protocol/accounting"

	mock "github.com/stretchr/testify/mock"
)

// MockIncomesHistory is an autogenerated mock type for the IncomesHistory type
type MockIncomesHistory struct {
	mock.Mock
}

type MockIncomesHistory_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIncomesHistory) EXPECT() *MockIncomesHistory_Expecter {
	return &MockIncomesHistory_Expecter{mock: &_m.Mock}
}

// GetHistory provides a mock function with given fields: ctx, params
func (_m *MockIncomesHistory) GetHistory(ctx context.Context, params *accounting.GetHistoryRequest) (*accounting.GetHistoryResponse, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for GetHistory")
	}

	var r0 *accounting.GetHistoryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *accounting.GetHistoryRequest) (*accounting.GetHistoryResponse, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *accounting.GetHistoryRequest) *accounting.GetHistoryResponse); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounting.GetHistoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *accounting.GetHistoryRequest) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIncomesHistory_GetHistory_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetHistory'
type MockIncomesHistory_GetHistory_Call struct {
	*mock.Call
}

// GetHistory is a helper method to define mock.On call
//   - ctx context.Context
//   - params *accounting.GetHistoryRequest
func (_e *MockIncomesHistory_Expecter) GetHistory(ctx interface{}, params interface{}) *MockIncomesHistory_GetHistory_Call {
	return &MockIncomesHistory_GetHistory_Call{Call: _e.mock.On("GetHistory", ctx, params)}
}

func (_c *MockIncomesHistory_GetHistory_Call) Run(run func(ctx context.Context, params *accounting.GetHistoryRequest)) *MockIncomesHistory_GetHistory_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*accounting.GetHistoryRequest))
	})
	return _c
}

func (_c *MockIncomesHistory_GetHistory_Call) Return(_a0 *accounting.GetHistoryResponse, _a1 error) *MockIncomesHistory_GetHistory_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIncomesHistory_GetHistory_Call) RunAndReturn(run func(context.Context, *accounting.GetHistoryRequest) (*accounting.GetHistoryResponse, error)) *MockIncomesHistory_GetHistory_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIncomesHistory creates a new instance of MockIncomesHistory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIncomesHistory(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIncomesHistory {
	mock := &MockIncomesHistory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Code generated by mockery v2.46.3. DO NOT EDIT.

package slack

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockSlack is an autogenerated mock type for the Slack type
type MockSlack struct {
	mock.Mock
}

type MockSlack_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSlack) EXPECT() *MockSlack_Expecter {
	return &MockSlack_Expecter{mock: &_m.Mock}
}

// Send provides a mock function with given fields: ctx, text
func (_m *MockSlack) Send(ctx context.Context, text string) error {
	ret := _m.Called(ctx, text)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, text)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSlack_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type MockSlack_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - ctx context.Context
//   - text string
func (_e *MockSlack_Expecter) Send(ctx interface{}, text interface{}) *MockSlack_Send_Call {
	return &MockSlack_Send_Call{Call: _e.mock.On("Send", ctx, text)}
}

func (_c *MockSlack_Send_Call) Run(run func(ctx context.Context, text string)) *MockSlack_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockSlack_Send_Call) Return(_a0 error) *MockSlack_Send_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSlack_Send_Call) RunAndReturn(run func(context.Context, string) error) *MockSlack_Send_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSlack creates a new instance of MockSlack. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSlack(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSlack {
	mock := &MockSlack{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
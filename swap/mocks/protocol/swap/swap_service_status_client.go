// Code generated by mockery v2.46.3. DO NOT EDIT.

package swap

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	metadata "google.golang.org/grpc/metadata"

	swap "code.emcdtech.com/b2b/swap/protocol/swap"
)

// MockSwapService_StatusClient is an autogenerated mock type for the SwapService_StatusClient type
type MockSwapService_StatusClient struct {
	mock.Mock
}

type MockSwapService_StatusClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSwapService_StatusClient) EXPECT() *MockSwapService_StatusClient_Expecter {
	return &MockSwapService_StatusClient_Expecter{mock: &_m.Mock}
}

// CloseSend provides a mock function with given fields:
func (_m *MockSwapService_StatusClient) CloseSend() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for CloseSend")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSwapService_StatusClient_CloseSend_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CloseSend'
type MockSwapService_StatusClient_CloseSend_Call struct {
	*mock.Call
}

// CloseSend is a helper method to define mock.On call
func (_e *MockSwapService_StatusClient_Expecter) CloseSend() *MockSwapService_StatusClient_CloseSend_Call {
	return &MockSwapService_StatusClient_CloseSend_Call{Call: _e.mock.On("CloseSend")}
}

func (_c *MockSwapService_StatusClient_CloseSend_Call) Run(run func()) *MockSwapService_StatusClient_CloseSend_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSwapService_StatusClient_CloseSend_Call) Return(_a0 error) *MockSwapService_StatusClient_CloseSend_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSwapService_StatusClient_CloseSend_Call) RunAndReturn(run func() error) *MockSwapService_StatusClient_CloseSend_Call {
	_c.Call.Return(run)
	return _c
}

// Context provides a mock function with given fields:
func (_m *MockSwapService_StatusClient) Context() context.Context {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Context")
	}

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// MockSwapService_StatusClient_Context_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Context'
type MockSwapService_StatusClient_Context_Call struct {
	*mock.Call
}

// Context is a helper method to define mock.On call
func (_e *MockSwapService_StatusClient_Expecter) Context() *MockSwapService_StatusClient_Context_Call {
	return &MockSwapService_StatusClient_Context_Call{Call: _e.mock.On("Context")}
}

func (_c *MockSwapService_StatusClient_Context_Call) Run(run func()) *MockSwapService_StatusClient_Context_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSwapService_StatusClient_Context_Call) Return(_a0 context.Context) *MockSwapService_StatusClient_Context_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSwapService_StatusClient_Context_Call) RunAndReturn(run func() context.Context) *MockSwapService_StatusClient_Context_Call {
	_c.Call.Return(run)
	return _c
}

// Header provides a mock function with given fields:
func (_m *MockSwapService_StatusClient) Header() (metadata.MD, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Header")
	}

	var r0 metadata.MD
	var r1 error
	if rf, ok := ret.Get(0).(func() (metadata.MD, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() metadata.MD); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(metadata.MD)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSwapService_StatusClient_Header_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Header'
type MockSwapService_StatusClient_Header_Call struct {
	*mock.Call
}

// Header is a helper method to define mock.On call
func (_e *MockSwapService_StatusClient_Expecter) Header() *MockSwapService_StatusClient_Header_Call {
	return &MockSwapService_StatusClient_Header_Call{Call: _e.mock.On("Header")}
}

func (_c *MockSwapService_StatusClient_Header_Call) Run(run func()) *MockSwapService_StatusClient_Header_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSwapService_StatusClient_Header_Call) Return(_a0 metadata.MD, _a1 error) *MockSwapService_StatusClient_Header_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSwapService_StatusClient_Header_Call) RunAndReturn(run func() (metadata.MD, error)) *MockSwapService_StatusClient_Header_Call {
	_c.Call.Return(run)
	return _c
}

// Recv provides a mock function with given fields:
func (_m *MockSwapService_StatusClient) Recv() (*swap.StatusResponse, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Recv")
	}

	var r0 *swap.StatusResponse
	var r1 error
	if rf, ok := ret.Get(0).(func() (*swap.StatusResponse, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *swap.StatusResponse); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*swap.StatusResponse)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSwapService_StatusClient_Recv_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Recv'
type MockSwapService_StatusClient_Recv_Call struct {
	*mock.Call
}

// Recv is a helper method to define mock.On call
func (_e *MockSwapService_StatusClient_Expecter) Recv() *MockSwapService_StatusClient_Recv_Call {
	return &MockSwapService_StatusClient_Recv_Call{Call: _e.mock.On("Recv")}
}

func (_c *MockSwapService_StatusClient_Recv_Call) Run(run func()) *MockSwapService_StatusClient_Recv_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSwapService_StatusClient_Recv_Call) Return(_a0 *swap.StatusResponse, _a1 error) *MockSwapService_StatusClient_Recv_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSwapService_StatusClient_Recv_Call) RunAndReturn(run func() (*swap.StatusResponse, error)) *MockSwapService_StatusClient_Recv_Call {
	_c.Call.Return(run)
	return _c
}

// RecvMsg provides a mock function with given fields: m
func (_m *MockSwapService_StatusClient) RecvMsg(m any) error {
	ret := _m.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for RecvMsg")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(any) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSwapService_StatusClient_RecvMsg_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RecvMsg'
type MockSwapService_StatusClient_RecvMsg_Call struct {
	*mock.Call
}

// RecvMsg is a helper method to define mock.On call
//   - m any
func (_e *MockSwapService_StatusClient_Expecter) RecvMsg(m interface{}) *MockSwapService_StatusClient_RecvMsg_Call {
	return &MockSwapService_StatusClient_RecvMsg_Call{Call: _e.mock.On("RecvMsg", m)}
}

func (_c *MockSwapService_StatusClient_RecvMsg_Call) Run(run func(m any)) *MockSwapService_StatusClient_RecvMsg_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(any))
	})
	return _c
}

func (_c *MockSwapService_StatusClient_RecvMsg_Call) Return(_a0 error) *MockSwapService_StatusClient_RecvMsg_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSwapService_StatusClient_RecvMsg_Call) RunAndReturn(run func(any) error) *MockSwapService_StatusClient_RecvMsg_Call {
	_c.Call.Return(run)
	return _c
}

// SendMsg provides a mock function with given fields: m
func (_m *MockSwapService_StatusClient) SendMsg(m any) error {
	ret := _m.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for SendMsg")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(any) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSwapService_StatusClient_SendMsg_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendMsg'
type MockSwapService_StatusClient_SendMsg_Call struct {
	*mock.Call
}

// SendMsg is a helper method to define mock.On call
//   - m any
func (_e *MockSwapService_StatusClient_Expecter) SendMsg(m interface{}) *MockSwapService_StatusClient_SendMsg_Call {
	return &MockSwapService_StatusClient_SendMsg_Call{Call: _e.mock.On("SendMsg", m)}
}

func (_c *MockSwapService_StatusClient_SendMsg_Call) Run(run func(m any)) *MockSwapService_StatusClient_SendMsg_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(any))
	})
	return _c
}

func (_c *MockSwapService_StatusClient_SendMsg_Call) Return(_a0 error) *MockSwapService_StatusClient_SendMsg_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSwapService_StatusClient_SendMsg_Call) RunAndReturn(run func(any) error) *MockSwapService_StatusClient_SendMsg_Call {
	_c.Call.Return(run)
	return _c
}

// Trailer provides a mock function with given fields:
func (_m *MockSwapService_StatusClient) Trailer() metadata.MD {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Trailer")
	}

	var r0 metadata.MD
	if rf, ok := ret.Get(0).(func() metadata.MD); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(metadata.MD)
		}
	}

	return r0
}

// MockSwapService_StatusClient_Trailer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Trailer'
type MockSwapService_StatusClient_Trailer_Call struct {
	*mock.Call
}

// Trailer is a helper method to define mock.On call
func (_e *MockSwapService_StatusClient_Expecter) Trailer() *MockSwapService_StatusClient_Trailer_Call {
	return &MockSwapService_StatusClient_Trailer_Call{Call: _e.mock.On("Trailer")}
}

func (_c *MockSwapService_StatusClient_Trailer_Call) Run(run func()) *MockSwapService_StatusClient_Trailer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSwapService_StatusClient_Trailer_Call) Return(_a0 metadata.MD) *MockSwapService_StatusClient_Trailer_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSwapService_StatusClient_Trailer_Call) RunAndReturn(run func() metadata.MD) *MockSwapService_StatusClient_Trailer_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSwapService_StatusClient creates a new instance of MockSwapService_StatusClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSwapService_StatusClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSwapService_StatusClient {
	mock := &MockSwapService_StatusClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

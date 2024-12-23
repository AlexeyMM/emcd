// Code generated by mockery v2.46.3. DO NOT EDIT.

package swapCoin

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	swapCoin "code.emcdtech.com/b2b/swap/protocol/swapCoin"
)

// MockSwapCoinServiceClient is an autogenerated mock type for the SwapCoinServiceClient type
type MockSwapCoinServiceClient struct {
	mock.Mock
}

type MockSwapCoinServiceClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSwapCoinServiceClient) EXPECT() *MockSwapCoinServiceClient_Expecter {
	return &MockSwapCoinServiceClient_Expecter{mock: &_m.Mock}
}

// GetAll provides a mock function with given fields: ctx, in, opts
func (_m *MockSwapCoinServiceClient) GetAll(ctx context.Context, in *swapCoin.GetAllRequest, opts ...grpc.CallOption) (*swapCoin.GetAllResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 *swapCoin.GetAllResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *swapCoin.GetAllRequest, ...grpc.CallOption) (*swapCoin.GetAllResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *swapCoin.GetAllRequest, ...grpc.CallOption) *swapCoin.GetAllResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*swapCoin.GetAllResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *swapCoin.GetAllRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSwapCoinServiceClient_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type MockSwapCoinServiceClient_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
//   - in *swapCoin.GetAllRequest
//   - opts ...grpc.CallOption
func (_e *MockSwapCoinServiceClient_Expecter) GetAll(ctx interface{}, in interface{}, opts ...interface{}) *MockSwapCoinServiceClient_GetAll_Call {
	return &MockSwapCoinServiceClient_GetAll_Call{Call: _e.mock.On("GetAll",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockSwapCoinServiceClient_GetAll_Call) Run(run func(ctx context.Context, in *swapCoin.GetAllRequest, opts ...grpc.CallOption)) *MockSwapCoinServiceClient_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*swapCoin.GetAllRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockSwapCoinServiceClient_GetAll_Call) Return(_a0 *swapCoin.GetAllResponse, _a1 error) *MockSwapCoinServiceClient_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSwapCoinServiceClient_GetAll_Call) RunAndReturn(run func(context.Context, *swapCoin.GetAllRequest, ...grpc.CallOption) (*swapCoin.GetAllResponse, error)) *MockSwapCoinServiceClient_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSwapCoinServiceClient creates a new instance of MockSwapCoinServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSwapCoinServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSwapCoinServiceClient {
	mock := &MockSwapCoinServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

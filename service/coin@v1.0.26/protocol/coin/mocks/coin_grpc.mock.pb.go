// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	coin "code.emcdtech.com/emcd/service/coin/protocol/coin"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// MockCoinServiceClient is an autogenerated mock type for the CoinServiceClient type
type MockCoinServiceClient struct {
	mock.Mock
}

type MockCoinServiceClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCoinServiceClient) EXPECT() *MockCoinServiceClient_Expecter {
	return &MockCoinServiceClient_Expecter{mock: &_m.Mock}
}

// GetCoin provides a mock function with given fields: ctx, in, opts
func (_m *MockCoinServiceClient) GetCoin(ctx context.Context, in *coin.GetCoinRequest, opts ...grpc.CallOption) (*coin.GetCoinResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetCoin")
	}

	var r0 *coin.GetCoinResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *coin.GetCoinRequest, ...grpc.CallOption) (*coin.GetCoinResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *coin.GetCoinRequest, ...grpc.CallOption) *coin.GetCoinResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coin.GetCoinResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *coin.GetCoinRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCoinServiceClient_GetCoin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCoin'
type MockCoinServiceClient_GetCoin_Call struct {
	*mock.Call
}

// GetCoin is a helper method to define mock.On call
//   - ctx context.Context
//   - in *coin.GetCoinRequest
//   - opts ...grpc.CallOption
func (_e *MockCoinServiceClient_Expecter) GetCoin(ctx interface{}, in interface{}, opts ...interface{}) *MockCoinServiceClient_GetCoin_Call {
	return &MockCoinServiceClient_GetCoin_Call{Call: _e.mock.On("GetCoin",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockCoinServiceClient_GetCoin_Call) Run(run func(ctx context.Context, in *coin.GetCoinRequest, opts ...grpc.CallOption)) *MockCoinServiceClient_GetCoin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*coin.GetCoinRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockCoinServiceClient_GetCoin_Call) Return(_a0 *coin.GetCoinResponse, _a1 error) *MockCoinServiceClient_GetCoin_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCoinServiceClient_GetCoin_Call) RunAndReturn(run func(context.Context, *coin.GetCoinRequest, ...grpc.CallOption) (*coin.GetCoinResponse, error)) *MockCoinServiceClient_GetCoin_Call {
	_c.Call.Return(run)
	return _c
}

// GetCoinIDFromLegacyID provides a mock function with given fields: ctx, in, opts
func (_m *MockCoinServiceClient) GetCoinIDFromLegacyID(ctx context.Context, in *coin.GetCoinIDFromLegacyIDRequest, opts ...grpc.CallOption) (*coin.GetCoinIDFromLegacyIDResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetCoinIDFromLegacyID")
	}

	var r0 *coin.GetCoinIDFromLegacyIDResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *coin.GetCoinIDFromLegacyIDRequest, ...grpc.CallOption) (*coin.GetCoinIDFromLegacyIDResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *coin.GetCoinIDFromLegacyIDRequest, ...grpc.CallOption) *coin.GetCoinIDFromLegacyIDResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coin.GetCoinIDFromLegacyIDResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *coin.GetCoinIDFromLegacyIDRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCoinServiceClient_GetCoinIDFromLegacyID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCoinIDFromLegacyID'
type MockCoinServiceClient_GetCoinIDFromLegacyID_Call struct {
	*mock.Call
}

// GetCoinIDFromLegacyID is a helper method to define mock.On call
//   - ctx context.Context
//   - in *coin.GetCoinIDFromLegacyIDRequest
//   - opts ...grpc.CallOption
func (_e *MockCoinServiceClient_Expecter) GetCoinIDFromLegacyID(ctx interface{}, in interface{}, opts ...interface{}) *MockCoinServiceClient_GetCoinIDFromLegacyID_Call {
	return &MockCoinServiceClient_GetCoinIDFromLegacyID_Call{Call: _e.mock.On("GetCoinIDFromLegacyID",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockCoinServiceClient_GetCoinIDFromLegacyID_Call) Run(run func(ctx context.Context, in *coin.GetCoinIDFromLegacyIDRequest, opts ...grpc.CallOption)) *MockCoinServiceClient_GetCoinIDFromLegacyID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*coin.GetCoinIDFromLegacyIDRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockCoinServiceClient_GetCoinIDFromLegacyID_Call) Return(_a0 *coin.GetCoinIDFromLegacyIDResponse, _a1 error) *MockCoinServiceClient_GetCoinIDFromLegacyID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCoinServiceClient_GetCoinIDFromLegacyID_Call) RunAndReturn(run func(context.Context, *coin.GetCoinIDFromLegacyIDRequest, ...grpc.CallOption) (*coin.GetCoinIDFromLegacyIDResponse, error)) *MockCoinServiceClient_GetCoinIDFromLegacyID_Call {
	_c.Call.Return(run)
	return _c
}

// GetCoins provides a mock function with given fields: ctx, in, opts
func (_m *MockCoinServiceClient) GetCoins(ctx context.Context, in *coin.GetCoinsRequest, opts ...grpc.CallOption) (*coin.GetCoinsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetCoins")
	}

	var r0 *coin.GetCoinsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *coin.GetCoinsRequest, ...grpc.CallOption) (*coin.GetCoinsResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *coin.GetCoinsRequest, ...grpc.CallOption) *coin.GetCoinsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coin.GetCoinsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *coin.GetCoinsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCoinServiceClient_GetCoins_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCoins'
type MockCoinServiceClient_GetCoins_Call struct {
	*mock.Call
}

// GetCoins is a helper method to define mock.On call
//   - ctx context.Context
//   - in *coin.GetCoinsRequest
//   - opts ...grpc.CallOption
func (_e *MockCoinServiceClient_Expecter) GetCoins(ctx interface{}, in interface{}, opts ...interface{}) *MockCoinServiceClient_GetCoins_Call {
	return &MockCoinServiceClient_GetCoins_Call{Call: _e.mock.On("GetCoins",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockCoinServiceClient_GetCoins_Call) Run(run func(ctx context.Context, in *coin.GetCoinsRequest, opts ...grpc.CallOption)) *MockCoinServiceClient_GetCoins_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*coin.GetCoinsRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockCoinServiceClient_GetCoins_Call) Return(_a0 *coin.GetCoinsResponse, _a1 error) *MockCoinServiceClient_GetCoins_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCoinServiceClient_GetCoins_Call) RunAndReturn(run func(context.Context, *coin.GetCoinsRequest, ...grpc.CallOption) (*coin.GetCoinsResponse, error)) *MockCoinServiceClient_GetCoins_Call {
	_c.Call.Return(run)
	return _c
}

// GetWithdrawalFee provides a mock function with given fields: ctx, in, opts
func (_m *MockCoinServiceClient) GetWithdrawalFee(ctx context.Context, in *coin.RequestGetWithdrawalFee, opts ...grpc.CallOption) (*coin.ResponseGetWithdrawalFee, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetWithdrawalFee")
	}

	var r0 *coin.ResponseGetWithdrawalFee
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *coin.RequestGetWithdrawalFee, ...grpc.CallOption) (*coin.ResponseGetWithdrawalFee, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *coin.RequestGetWithdrawalFee, ...grpc.CallOption) *coin.ResponseGetWithdrawalFee); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coin.ResponseGetWithdrawalFee)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *coin.RequestGetWithdrawalFee, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCoinServiceClient_GetWithdrawalFee_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWithdrawalFee'
type MockCoinServiceClient_GetWithdrawalFee_Call struct {
	*mock.Call
}

// GetWithdrawalFee is a helper method to define mock.On call
//   - ctx context.Context
//   - in *coin.RequestGetWithdrawalFee
//   - opts ...grpc.CallOption
func (_e *MockCoinServiceClient_Expecter) GetWithdrawalFee(ctx interface{}, in interface{}, opts ...interface{}) *MockCoinServiceClient_GetWithdrawalFee_Call {
	return &MockCoinServiceClient_GetWithdrawalFee_Call{Call: _e.mock.On("GetWithdrawalFee",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockCoinServiceClient_GetWithdrawalFee_Call) Run(run func(ctx context.Context, in *coin.RequestGetWithdrawalFee, opts ...grpc.CallOption)) *MockCoinServiceClient_GetWithdrawalFee_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*coin.RequestGetWithdrawalFee), variadicArgs...)
	})
	return _c
}

func (_c *MockCoinServiceClient_GetWithdrawalFee_Call) Return(_a0 *coin.ResponseGetWithdrawalFee, _a1 error) *MockCoinServiceClient_GetWithdrawalFee_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCoinServiceClient_GetWithdrawalFee_Call) RunAndReturn(run func(context.Context, *coin.RequestGetWithdrawalFee, ...grpc.CallOption) (*coin.ResponseGetWithdrawalFee, error)) *MockCoinServiceClient_GetWithdrawalFee_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCoinServiceClient creates a new instance of MockCoinServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCoinServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCoinServiceClient {
	mock := &MockCoinServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Code generated by mockery v2.46.3. DO NOT EDIT.

package repository

import (
	context "context"

	model "code.emcdtech.com/b2b/swap/model"
	mock "github.com/stretchr/testify/mock"
)

// MockCoin is an autogenerated mock type for the Coin type
type MockCoin struct {
	mock.Mock
}

type MockCoin_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCoin) EXPECT() *MockCoin_Expecter {
	return &MockCoin_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: ctx, coin
func (_m *MockCoin) Get(ctx context.Context, coin string) (*model.Coin, error) {
	ret := _m.Called(ctx, coin)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *model.Coin
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Coin, error)); ok {
		return rf(ctx, coin)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Coin); ok {
		r0 = rf(ctx, coin)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Coin)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, coin)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCoin_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockCoin_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - coin string
func (_e *MockCoin_Expecter) Get(ctx interface{}, coin interface{}) *MockCoin_Get_Call {
	return &MockCoin_Get_Call{Call: _e.mock.On("Get", ctx, coin)}
}

func (_c *MockCoin_Get_Call) Run(run func(ctx context.Context, coin string)) *MockCoin_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockCoin_Get_Call) Return(_a0 *model.Coin, _a1 error) *MockCoin_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCoin_Get_Call) RunAndReturn(run func(context.Context, string) (*model.Coin, error)) *MockCoin_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetAccuracyForWithdrawAndDeposit provides a mock function with given fields: ctx, coin, network
func (_m *MockCoin) GetAccuracyForWithdrawAndDeposit(ctx context.Context, coin string, network string) (int, error) {
	ret := _m.Called(ctx, coin, network)

	if len(ret) == 0 {
		panic("no return value specified for GetAccuracyForWithdrawAndDeposit")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (int, error)); ok {
		return rf(ctx, coin, network)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) int); ok {
		r0 = rf(ctx, coin, network)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, coin, network)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCoin_GetAccuracyForWithdrawAndDeposit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAccuracyForWithdrawAndDeposit'
type MockCoin_GetAccuracyForWithdrawAndDeposit_Call struct {
	*mock.Call
}

// GetAccuracyForWithdrawAndDeposit is a helper method to define mock.On call
//   - ctx context.Context
//   - coin string
//   - network string
func (_e *MockCoin_Expecter) GetAccuracyForWithdrawAndDeposit(ctx interface{}, coin interface{}, network interface{}) *MockCoin_GetAccuracyForWithdrawAndDeposit_Call {
	return &MockCoin_GetAccuracyForWithdrawAndDeposit_Call{Call: _e.mock.On("GetAccuracyForWithdrawAndDeposit", ctx, coin, network)}
}

func (_c *MockCoin_GetAccuracyForWithdrawAndDeposit_Call) Run(run func(ctx context.Context, coin string, network string)) *MockCoin_GetAccuracyForWithdrawAndDeposit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockCoin_GetAccuracyForWithdrawAndDeposit_Call) Return(_a0 int, _a1 error) *MockCoin_GetAccuracyForWithdrawAndDeposit_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCoin_GetAccuracyForWithdrawAndDeposit_Call) RunAndReturn(run func(context.Context, string, string) (int, error)) *MockCoin_GetAccuracyForWithdrawAndDeposit_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields: ctx
func (_m *MockCoin) GetAll(ctx context.Context) ([]*model.Coin, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []*model.Coin
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*model.Coin, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*model.Coin); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Coin)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCoin_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type MockCoin_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockCoin_Expecter) GetAll(ctx interface{}) *MockCoin_GetAll_Call {
	return &MockCoin_GetAll_Call{Call: _e.mock.On("GetAll", ctx)}
}

func (_c *MockCoin_GetAll_Call) Run(run func(ctx context.Context)) *MockCoin_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockCoin_GetAll_Call) Return(_a0 []*model.Coin, _a1 error) *MockCoin_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCoin_GetAll_Call) RunAndReturn(run func(context.Context) ([]*model.Coin, error)) *MockCoin_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetAnyNetwork provides a mock function with given fields: ctx, coin
func (_m *MockCoin) GetAnyNetwork(ctx context.Context, coin string) (*model.Network, error) {
	ret := _m.Called(ctx, coin)

	if len(ret) == 0 {
		panic("no return value specified for GetAnyNetwork")
	}

	var r0 *model.Network
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Network, error)); ok {
		return rf(ctx, coin)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Network); ok {
		r0 = rf(ctx, coin)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Network)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, coin)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCoin_GetAnyNetwork_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAnyNetwork'
type MockCoin_GetAnyNetwork_Call struct {
	*mock.Call
}

// GetAnyNetwork is a helper method to define mock.On call
//   - ctx context.Context
//   - coin string
func (_e *MockCoin_Expecter) GetAnyNetwork(ctx interface{}, coin interface{}) *MockCoin_GetAnyNetwork_Call {
	return &MockCoin_GetAnyNetwork_Call{Call: _e.mock.On("GetAnyNetwork", ctx, coin)}
}

func (_c *MockCoin_GetAnyNetwork_Call) Run(run func(ctx context.Context, coin string)) *MockCoin_GetAnyNetwork_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockCoin_GetAnyNetwork_Call) Return(_a0 *model.Network, _a1 error) *MockCoin_GetAnyNetwork_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCoin_GetAnyNetwork_Call) RunAndReturn(run func(context.Context, string) (*model.Network, error)) *MockCoin_GetAnyNetwork_Call {
	_c.Call.Return(run)
	return _c
}

// GetNetwork provides a mock function with given fields: ctx, coin, network
func (_m *MockCoin) GetNetwork(ctx context.Context, coin string, network string) (*model.Network, error) {
	ret := _m.Called(ctx, coin, network)

	if len(ret) == 0 {
		panic("no return value specified for GetNetwork")
	}

	var r0 *model.Network
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*model.Network, error)); ok {
		return rf(ctx, coin, network)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.Network); ok {
		r0 = rf(ctx, coin, network)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Network)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, coin, network)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCoin_GetNetwork_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNetwork'
type MockCoin_GetNetwork_Call struct {
	*mock.Call
}

// GetNetwork is a helper method to define mock.On call
//   - ctx context.Context
//   - coin string
//   - network string
func (_e *MockCoin_Expecter) GetNetwork(ctx interface{}, coin interface{}, network interface{}) *MockCoin_GetNetwork_Call {
	return &MockCoin_GetNetwork_Call{Call: _e.mock.On("GetNetwork", ctx, coin, network)}
}

func (_c *MockCoin_GetNetwork_Call) Run(run func(ctx context.Context, coin string, network string)) *MockCoin_GetNetwork_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockCoin_GetNetwork_Call) Return(_a0 *model.Network, _a1 error) *MockCoin_GetNetwork_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCoin_GetNetwork_Call) RunAndReturn(run func(context.Context, string, string) (*model.Network, error)) *MockCoin_GetNetwork_Call {
	_c.Call.Return(run)
	return _c
}

// GetWithdrawFee provides a mock function with given fields: ctx, coin, network
func (_m *MockCoin) GetWithdrawFee(ctx context.Context, coin string, network string) (*model.WithdrawFee, error) {
	ret := _m.Called(ctx, coin, network)

	if len(ret) == 0 {
		panic("no return value specified for GetWithdrawFee")
	}

	var r0 *model.WithdrawFee
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*model.WithdrawFee, error)); ok {
		return rf(ctx, coin, network)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.WithdrawFee); ok {
		r0 = rf(ctx, coin, network)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.WithdrawFee)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, coin, network)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCoin_GetWithdrawFee_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWithdrawFee'
type MockCoin_GetWithdrawFee_Call struct {
	*mock.Call
}

// GetWithdrawFee is a helper method to define mock.On call
//   - ctx context.Context
//   - coin string
//   - network string
func (_e *MockCoin_Expecter) GetWithdrawFee(ctx interface{}, coin interface{}, network interface{}) *MockCoin_GetWithdrawFee_Call {
	return &MockCoin_GetWithdrawFee_Call{Call: _e.mock.On("GetWithdrawFee", ctx, coin, network)}
}

func (_c *MockCoin_GetWithdrawFee_Call) Run(run func(ctx context.Context, coin string, network string)) *MockCoin_GetWithdrawFee_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockCoin_GetWithdrawFee_Call) Return(_a0 *model.WithdrawFee, _a1 error) *MockCoin_GetWithdrawFee_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCoin_GetWithdrawFee_Call) RunAndReturn(run func(context.Context, string, string) (*model.WithdrawFee, error)) *MockCoin_GetWithdrawFee_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateAll provides a mock function with given fields: ctx, coins
func (_m *MockCoin) UpdateAll(ctx context.Context, coins []*model.Coin) error {
	ret := _m.Called(ctx, coins)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAll")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*model.Coin) error); ok {
		r0 = rf(ctx, coins)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCoin_UpdateAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateAll'
type MockCoin_UpdateAll_Call struct {
	*mock.Call
}

// UpdateAll is a helper method to define mock.On call
//   - ctx context.Context
//   - coins []*model.Coin
func (_e *MockCoin_Expecter) UpdateAll(ctx interface{}, coins interface{}) *MockCoin_UpdateAll_Call {
	return &MockCoin_UpdateAll_Call{Call: _e.mock.On("UpdateAll", ctx, coins)}
}

func (_c *MockCoin_UpdateAll_Call) Run(run func(ctx context.Context, coins []*model.Coin)) *MockCoin_UpdateAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*model.Coin))
	})
	return _c
}

func (_c *MockCoin_UpdateAll_Call) Return(_a0 error) *MockCoin_UpdateAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCoin_UpdateAll_Call) RunAndReturn(run func(context.Context, []*model.Coin) error) *MockCoin_UpdateAll_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCoin creates a new instance of MockCoin. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCoin(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCoin {
	mock := &MockCoin{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

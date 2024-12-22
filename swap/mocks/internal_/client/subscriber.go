// Code generated by mockery v2.46.3. DO NOT EDIT.

package client

import (
	context "context"

	decimal "github.com/shopspring/decimal"
	mock "github.com/stretchr/testify/mock"

	model "code.emcdtech.com/b2b/swap/model"

	time "time"

	uuid "github.com/google/uuid"
)

// MockSubscriber is an autogenerated mock type for the Subscriber type
type MockSubscriber struct {
	mock.Mock
}

type MockSubscriber_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSubscriber) EXPECT() *MockSubscriber_Expecter {
	return &MockSubscriber_Expecter{mock: &_m.Mock}
}

// SubscribeOnOrderbooks provides a mock function with given fields: ctx, symbols
func (_m *MockSubscriber) SubscribeOnOrderbooks(ctx context.Context, symbols []*model.Symbol) error {
	ret := _m.Called(ctx, symbols)

	if len(ret) == 0 {
		panic("no return value specified for SubscribeOnOrderbooks")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*model.Symbol) error); ok {
		r0 = rf(ctx, symbols)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSubscriber_SubscribeOnOrderbooks_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SubscribeOnOrderbooks'
type MockSubscriber_SubscribeOnOrderbooks_Call struct {
	*mock.Call
}

// SubscribeOnOrderbooks is a helper method to define mock.On call
//   - ctx context.Context
//   - symbols []*model.Symbol
func (_e *MockSubscriber_Expecter) SubscribeOnOrderbooks(ctx interface{}, symbols interface{}) *MockSubscriber_SubscribeOnOrderbooks_Call {
	return &MockSubscriber_SubscribeOnOrderbooks_Call{Call: _e.mock.On("SubscribeOnOrderbooks", ctx, symbols)}
}

func (_c *MockSubscriber_SubscribeOnOrderbooks_Call) Run(run func(ctx context.Context, symbols []*model.Symbol)) *MockSubscriber_SubscribeOnOrderbooks_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*model.Symbol))
	})
	return _c
}

func (_c *MockSubscriber_SubscribeOnOrderbooks_Call) Return(_a0 error) *MockSubscriber_SubscribeOnOrderbooks_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSubscriber_SubscribeOnOrderbooks_Call) RunAndReturn(run func(context.Context, []*model.Symbol) error) *MockSubscriber_SubscribeOnOrderbooks_Call {
	_c.Call.Return(run)
	return _c
}

// SubscribeOnOrders provides a mock function with given fields: ctx, account, orders, receivedFirstOrder
func (_m *MockSubscriber) SubscribeOnOrders(ctx context.Context, account *model.Account, orders []*model.Order, receivedFirstOrder bool) error {
	ret := _m.Called(ctx, account, orders, receivedFirstOrder)

	if len(ret) == 0 {
		panic("no return value specified for SubscribeOnOrders")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Account, []*model.Order, bool) error); ok {
		r0 = rf(ctx, account, orders, receivedFirstOrder)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSubscriber_SubscribeOnOrders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SubscribeOnOrders'
type MockSubscriber_SubscribeOnOrders_Call struct {
	*mock.Call
}

// SubscribeOnOrders is a helper method to define mock.On call
//   - ctx context.Context
//   - account *model.Account
//   - orders []*model.Order
//   - receivedFirstOrder bool
func (_e *MockSubscriber_Expecter) SubscribeOnOrders(ctx interface{}, account interface{}, orders interface{}, receivedFirstOrder interface{}) *MockSubscriber_SubscribeOnOrders_Call {
	return &MockSubscriber_SubscribeOnOrders_Call{Call: _e.mock.On("SubscribeOnOrders", ctx, account, orders, receivedFirstOrder)}
}

func (_c *MockSubscriber_SubscribeOnOrders_Call) Run(run func(ctx context.Context, account *model.Account, orders []*model.Order, receivedFirstOrder bool)) *MockSubscriber_SubscribeOnOrders_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Account), args[2].([]*model.Order), args[3].(bool))
	})
	return _c
}

func (_c *MockSubscriber_SubscribeOnOrders_Call) Return(_a0 error) *MockSubscriber_SubscribeOnOrders_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSubscriber_SubscribeOnOrders_Call) RunAndReturn(run func(context.Context, *model.Account, []*model.Order, bool) error) *MockSubscriber_SubscribeOnOrders_Call {
	_c.Call.Return(run)
	return _c
}

// SubscribeOnWallet provides a mock function with given fields: ctx, swapID, swapCreated, account, coin, amount
func (_m *MockSubscriber) SubscribeOnWallet(ctx context.Context, swapID uuid.UUID, swapCreated time.Time, account *model.Account, coin string, amount decimal.Decimal) error {
	ret := _m.Called(ctx, swapID, swapCreated, account, coin, amount)

	if len(ret) == 0 {
		panic("no return value specified for SubscribeOnWallet")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, time.Time, *model.Account, string, decimal.Decimal) error); ok {
		r0 = rf(ctx, swapID, swapCreated, account, coin, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSubscriber_SubscribeOnWallet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SubscribeOnWallet'
type MockSubscriber_SubscribeOnWallet_Call struct {
	*mock.Call
}

// SubscribeOnWallet is a helper method to define mock.On call
//   - ctx context.Context
//   - swapID uuid.UUID
//   - swapCreated time.Time
//   - account *model.Account
//   - coin string
//   - amount decimal.Decimal
func (_e *MockSubscriber_Expecter) SubscribeOnWallet(ctx interface{}, swapID interface{}, swapCreated interface{}, account interface{}, coin interface{}, amount interface{}) *MockSubscriber_SubscribeOnWallet_Call {
	return &MockSubscriber_SubscribeOnWallet_Call{Call: _e.mock.On("SubscribeOnWallet", ctx, swapID, swapCreated, account, coin, amount)}
}

func (_c *MockSubscriber_SubscribeOnWallet_Call) Run(run func(ctx context.Context, swapID uuid.UUID, swapCreated time.Time, account *model.Account, coin string, amount decimal.Decimal)) *MockSubscriber_SubscribeOnWallet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(time.Time), args[3].(*model.Account), args[4].(string), args[5].(decimal.Decimal))
	})
	return _c
}

func (_c *MockSubscriber_SubscribeOnWallet_Call) Return(_a0 error) *MockSubscriber_SubscribeOnWallet_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSubscriber_SubscribeOnWallet_Call) RunAndReturn(run func(context.Context, uuid.UUID, time.Time, *model.Account, string, decimal.Decimal) error) *MockSubscriber_SubscribeOnWallet_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSubscriber creates a new instance of MockSubscriber. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSubscriber(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSubscriber {
	mock := &MockSubscriber{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

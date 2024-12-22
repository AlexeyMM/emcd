// Code generated by mockery v2.46.3. DO NOT EDIT.

package service

import (
	context "context"

	decimal "github.com/shopspring/decimal"
	mock "github.com/stretchr/testify/mock"

	model "code.emcdtech.com/b2b/swap/model"

	uuid "github.com/google/uuid"
)

// MockAdmin is an autogenerated mock type for the Admin type
type MockAdmin struct {
	mock.Mock
}

type MockAdmin_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAdmin) EXPECT() *MockAdmin_Expecter {
	return &MockAdmin_Expecter{mock: &_m.Mock}
}

// ChangeManualSwapStatus provides a mock function with given fields: ctx, swapID, status
func (_m *MockAdmin) ChangeManualSwapStatus(ctx context.Context, swapID uuid.UUID, status model.Status) error {
	ret := _m.Called(ctx, swapID, status)

	if len(ret) == 0 {
		panic("no return value specified for ChangeManualSwapStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, model.Status) error); ok {
		r0 = rf(ctx, swapID, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAdmin_ChangeManualSwapStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ChangeManualSwapStatus'
type MockAdmin_ChangeManualSwapStatus_Call struct {
	*mock.Call
}

// ChangeManualSwapStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - swapID uuid.UUID
//   - status model.Status
func (_e *MockAdmin_Expecter) ChangeManualSwapStatus(ctx interface{}, swapID interface{}, status interface{}) *MockAdmin_ChangeManualSwapStatus_Call {
	return &MockAdmin_ChangeManualSwapStatus_Call{Call: _e.mock.On("ChangeManualSwapStatus", ctx, swapID, status)}
}

func (_c *MockAdmin_ChangeManualSwapStatus_Call) Run(run func(ctx context.Context, swapID uuid.UUID, status model.Status)) *MockAdmin_ChangeManualSwapStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(model.Status))
	})
	return _c
}

func (_c *MockAdmin_ChangeManualSwapStatus_Call) Return(_a0 error) *MockAdmin_ChangeManualSwapStatus_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAdmin_ChangeManualSwapStatus_Call) RunAndReturn(run func(context.Context, uuid.UUID, model.Status) error) *MockAdmin_ChangeManualSwapStatus_Call {
	_c.Call.Return(run)
	return _c
}

// CheckOrder provides a mock function with given fields: ctx, id
func (_m *MockAdmin) CheckOrder(ctx context.Context, id uuid.UUID) (model.OrderStatus, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for CheckOrder")
	}

	var r0 model.OrderStatus
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (model.OrderStatus, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) model.OrderStatus); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(model.OrderStatus)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAdmin_CheckOrder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckOrder'
type MockAdmin_CheckOrder_Call struct {
	*mock.Call
}

// CheckOrder is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
func (_e *MockAdmin_Expecter) CheckOrder(ctx interface{}, id interface{}) *MockAdmin_CheckOrder_Call {
	return &MockAdmin_CheckOrder_Call{Call: _e.mock.On("CheckOrder", ctx, id)}
}

func (_c *MockAdmin_CheckOrder_Call) Run(run func(ctx context.Context, id uuid.UUID)) *MockAdmin_CheckOrder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockAdmin_CheckOrder_Call) Return(_a0 model.OrderStatus, _a1 error) *MockAdmin_CheckOrder_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAdmin_CheckOrder_Call) RunAndReturn(run func(context.Context, uuid.UUID) (model.OrderStatus, error)) *MockAdmin_CheckOrder_Call {
	_c.Call.Return(run)
	return _c
}

// ConfirmAQuote provides a mock function with given fields: ctx, id
func (_m *MockAdmin) ConfirmAQuote(ctx context.Context, id string) (string, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for ConfirmAQuote")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAdmin_ConfirmAQuote_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConfirmAQuote'
type MockAdmin_ConfirmAQuote_Call struct {
	*mock.Call
}

// ConfirmAQuote is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *MockAdmin_Expecter) ConfirmAQuote(ctx interface{}, id interface{}) *MockAdmin_ConfirmAQuote_Call {
	return &MockAdmin_ConfirmAQuote_Call{Call: _e.mock.On("ConfirmAQuote", ctx, id)}
}

func (_c *MockAdmin_ConfirmAQuote_Call) Run(run func(ctx context.Context, id string)) *MockAdmin_ConfirmAQuote_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockAdmin_ConfirmAQuote_Call) Return(_a0 string, _a1 error) *MockAdmin_ConfirmAQuote_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAdmin_ConfirmAQuote_Call) RunAndReturn(run func(context.Context, string) (string, error)) *MockAdmin_ConfirmAQuote_Call {
	_c.Call.Return(run)
	return _c
}

// GetBalanceByCoin provides a mock function with given fields: ctx, accountType, coin
func (_m *MockAdmin) GetBalanceByCoin(ctx context.Context, accountType string, coin string) (decimal.Decimal, error) {
	ret := _m.Called(ctx, accountType, coin)

	if len(ret) == 0 {
		panic("no return value specified for GetBalanceByCoin")
	}

	var r0 decimal.Decimal
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (decimal.Decimal, error)); ok {
		return rf(ctx, accountType, coin)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) decimal.Decimal); ok {
		r0 = rf(ctx, accountType, coin)
	} else {
		r0 = ret.Get(0).(decimal.Decimal)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, accountType, coin)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAdmin_GetBalanceByCoin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBalanceByCoin'
type MockAdmin_GetBalanceByCoin_Call struct {
	*mock.Call
}

// GetBalanceByCoin is a helper method to define mock.On call
//   - ctx context.Context
//   - accountType string
//   - coin string
func (_e *MockAdmin_Expecter) GetBalanceByCoin(ctx interface{}, accountType interface{}, coin interface{}) *MockAdmin_GetBalanceByCoin_Call {
	return &MockAdmin_GetBalanceByCoin_Call{Call: _e.mock.On("GetBalanceByCoin", ctx, accountType, coin)}
}

func (_c *MockAdmin_GetBalanceByCoin_Call) Run(run func(ctx context.Context, accountType string, coin string)) *MockAdmin_GetBalanceByCoin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockAdmin_GetBalanceByCoin_Call) Return(_a0 decimal.Decimal, _a1 error) *MockAdmin_GetBalanceByCoin_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAdmin_GetBalanceByCoin_Call) RunAndReturn(run func(context.Context, string, string) (decimal.Decimal, error)) *MockAdmin_GetBalanceByCoin_Call {
	_c.Call.Return(run)
	return _c
}

// GetConvertStatus provides a mock function with given fields: ctx, id, accountType
func (_m *MockAdmin) GetConvertStatus(ctx context.Context, id string, accountType string) (string, error) {
	ret := _m.Called(ctx, id, accountType)

	if len(ret) == 0 {
		panic("no return value specified for GetConvertStatus")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (string, error)); ok {
		return rf(ctx, id, accountType)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, id, accountType)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, id, accountType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAdmin_GetConvertStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConvertStatus'
type MockAdmin_GetConvertStatus_Call struct {
	*mock.Call
}

// GetConvertStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
//   - accountType string
func (_e *MockAdmin_Expecter) GetConvertStatus(ctx interface{}, id interface{}, accountType interface{}) *MockAdmin_GetConvertStatus_Call {
	return &MockAdmin_GetConvertStatus_Call{Call: _e.mock.On("GetConvertStatus", ctx, id, accountType)}
}

func (_c *MockAdmin_GetConvertStatus_Call) Run(run func(ctx context.Context, id string, accountType string)) *MockAdmin_GetConvertStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockAdmin_GetConvertStatus_Call) Return(_a0 string, _a1 error) *MockAdmin_GetConvertStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAdmin_GetConvertStatus_Call) RunAndReturn(run func(context.Context, string, string) (string, error)) *MockAdmin_GetConvertStatus_Call {
	_c.Call.Return(run)
	return _c
}

// GetSwapStatusHistory provides a mock function with given fields: ctx, swapID
func (_m *MockAdmin) GetSwapStatusHistory(ctx context.Context, swapID uuid.UUID) ([]*model.SwapStatusHistoryItem, error) {
	ret := _m.Called(ctx, swapID)

	if len(ret) == 0 {
		panic("no return value specified for GetSwapStatusHistory")
	}

	var r0 []*model.SwapStatusHistoryItem
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]*model.SwapStatusHistoryItem, error)); ok {
		return rf(ctx, swapID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []*model.SwapStatusHistoryItem); ok {
		r0 = rf(ctx, swapID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.SwapStatusHistoryItem)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, swapID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAdmin_GetSwapStatusHistory_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSwapStatusHistory'
type MockAdmin_GetSwapStatusHistory_Call struct {
	*mock.Call
}

// GetSwapStatusHistory is a helper method to define mock.On call
//   - ctx context.Context
//   - swapID uuid.UUID
func (_e *MockAdmin_Expecter) GetSwapStatusHistory(ctx interface{}, swapID interface{}) *MockAdmin_GetSwapStatusHistory_Call {
	return &MockAdmin_GetSwapStatusHistory_Call{Call: _e.mock.On("GetSwapStatusHistory", ctx, swapID)}
}

func (_c *MockAdmin_GetSwapStatusHistory_Call) Run(run func(ctx context.Context, swapID uuid.UUID)) *MockAdmin_GetSwapStatusHistory_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockAdmin_GetSwapStatusHistory_Call) Return(_a0 []*model.SwapStatusHistoryItem, _a1 error) *MockAdmin_GetSwapStatusHistory_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAdmin_GetSwapStatusHistory_Call) RunAndReturn(run func(context.Context, uuid.UUID) ([]*model.SwapStatusHistoryItem, error)) *MockAdmin_GetSwapStatusHistory_Call {
	_c.Call.Return(run)
	return _c
}

// GetWithdrawalLink provides a mock function with given fields: ctx, withdrawalID
func (_m *MockAdmin) GetWithdrawalLink(ctx context.Context, withdrawalID int) (string, error) {
	ret := _m.Called(ctx, withdrawalID)

	if len(ret) == 0 {
		panic("no return value specified for GetWithdrawalLink")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (string, error)); ok {
		return rf(ctx, withdrawalID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) string); ok {
		r0 = rf(ctx, withdrawalID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, withdrawalID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAdmin_GetWithdrawalLink_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWithdrawalLink'
type MockAdmin_GetWithdrawalLink_Call struct {
	*mock.Call
}

// GetWithdrawalLink is a helper method to define mock.On call
//   - ctx context.Context
//   - withdrawalID int
func (_e *MockAdmin_Expecter) GetWithdrawalLink(ctx interface{}, withdrawalID interface{}) *MockAdmin_GetWithdrawalLink_Call {
	return &MockAdmin_GetWithdrawalLink_Call{Call: _e.mock.On("GetWithdrawalLink", ctx, withdrawalID)}
}

func (_c *MockAdmin_GetWithdrawalLink_Call) Run(run func(ctx context.Context, withdrawalID int)) *MockAdmin_GetWithdrawalLink_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *MockAdmin_GetWithdrawalLink_Call) Return(_a0 string, _a1 error) *MockAdmin_GetWithdrawalLink_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAdmin_GetWithdrawalLink_Call) RunAndReturn(run func(context.Context, int) (string, error)) *MockAdmin_GetWithdrawalLink_Call {
	_c.Call.Return(run)
	return _c
}

// PlaceOrderForUSDT provides a mock function with given fields: ctx, coin, direction, amount
func (_m *MockAdmin) PlaceOrderForUSDT(ctx context.Context, coin string, direction model.Direction, amount decimal.Decimal) (uuid.UUID, error) {
	ret := _m.Called(ctx, coin, direction, amount)

	if len(ret) == 0 {
		panic("no return value specified for PlaceOrderForUSDT")
	}

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.Direction, decimal.Decimal) (uuid.UUID, error)); ok {
		return rf(ctx, coin, direction, amount)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, model.Direction, decimal.Decimal) uuid.UUID); ok {
		r0 = rf(ctx, coin, direction, amount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, model.Direction, decimal.Decimal) error); ok {
		r1 = rf(ctx, coin, direction, amount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAdmin_PlaceOrderForUSDT_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PlaceOrderForUSDT'
type MockAdmin_PlaceOrderForUSDT_Call struct {
	*mock.Call
}

// PlaceOrderForUSDT is a helper method to define mock.On call
//   - ctx context.Context
//   - coin string
//   - direction model.Direction
//   - amount decimal.Decimal
func (_e *MockAdmin_Expecter) PlaceOrderForUSDT(ctx interface{}, coin interface{}, direction interface{}, amount interface{}) *MockAdmin_PlaceOrderForUSDT_Call {
	return &MockAdmin_PlaceOrderForUSDT_Call{Call: _e.mock.On("PlaceOrderForUSDT", ctx, coin, direction, amount)}
}

func (_c *MockAdmin_PlaceOrderForUSDT_Call) Run(run func(ctx context.Context, coin string, direction model.Direction, amount decimal.Decimal)) *MockAdmin_PlaceOrderForUSDT_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(model.Direction), args[3].(decimal.Decimal))
	})
	return _c
}

func (_c *MockAdmin_PlaceOrderForUSDT_Call) Return(_a0 uuid.UUID, _a1 error) *MockAdmin_PlaceOrderForUSDT_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAdmin_PlaceOrderForUSDT_Call) RunAndReturn(run func(context.Context, string, model.Direction, decimal.Decimal) (uuid.UUID, error)) *MockAdmin_PlaceOrderForUSDT_Call {
	_c.Call.Return(run)
	return _c
}

// RequestAQuote provides a mock function with given fields: ctx, from, to, accountType, amount
func (_m *MockAdmin) RequestAQuote(ctx context.Context, from string, to string, accountType string, amount decimal.Decimal) (*model.Quote, error) {
	ret := _m.Called(ctx, from, to, accountType, amount)

	if len(ret) == 0 {
		panic("no return value specified for RequestAQuote")
	}

	var r0 *model.Quote
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, decimal.Decimal) (*model.Quote, error)); ok {
		return rf(ctx, from, to, accountType, amount)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, decimal.Decimal) *model.Quote); ok {
		r0 = rf(ctx, from, to, accountType, amount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Quote)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, decimal.Decimal) error); ok {
		r1 = rf(ctx, from, to, accountType, amount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAdmin_RequestAQuote_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RequestAQuote'
type MockAdmin_RequestAQuote_Call struct {
	*mock.Call
}

// RequestAQuote is a helper method to define mock.On call
//   - ctx context.Context
//   - from string
//   - to string
//   - accountType string
//   - amount decimal.Decimal
func (_e *MockAdmin_Expecter) RequestAQuote(ctx interface{}, from interface{}, to interface{}, accountType interface{}, amount interface{}) *MockAdmin_RequestAQuote_Call {
	return &MockAdmin_RequestAQuote_Call{Call: _e.mock.On("RequestAQuote", ctx, from, to, accountType, amount)}
}

func (_c *MockAdmin_RequestAQuote_Call) Run(run func(ctx context.Context, from string, to string, accountType string, amount decimal.Decimal)) *MockAdmin_RequestAQuote_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string), args[4].(decimal.Decimal))
	})
	return _c
}

func (_c *MockAdmin_RequestAQuote_Call) Return(_a0 *model.Quote, _a1 error) *MockAdmin_RequestAQuote_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAdmin_RequestAQuote_Call) RunAndReturn(run func(context.Context, string, string, string, decimal.Decimal) (*model.Quote, error)) *MockAdmin_RequestAQuote_Call {
	_c.Call.Return(run)
	return _c
}

// TransferBetweenAccountTypes provides a mock function with given fields: ctx, fromAccountType, ToAccountType, coin, amount
func (_m *MockAdmin) TransferBetweenAccountTypes(ctx context.Context, fromAccountType string, ToAccountType string, coin string, amount decimal.Decimal) error {
	ret := _m.Called(ctx, fromAccountType, ToAccountType, coin, amount)

	if len(ret) == 0 {
		panic("no return value specified for TransferBetweenAccountTypes")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, decimal.Decimal) error); ok {
		r0 = rf(ctx, fromAccountType, ToAccountType, coin, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAdmin_TransferBetweenAccountTypes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TransferBetweenAccountTypes'
type MockAdmin_TransferBetweenAccountTypes_Call struct {
	*mock.Call
}

// TransferBetweenAccountTypes is a helper method to define mock.On call
//   - ctx context.Context
//   - fromAccountType string
//   - ToAccountType string
//   - coin string
//   - amount decimal.Decimal
func (_e *MockAdmin_Expecter) TransferBetweenAccountTypes(ctx interface{}, fromAccountType interface{}, ToAccountType interface{}, coin interface{}, amount interface{}) *MockAdmin_TransferBetweenAccountTypes_Call {
	return &MockAdmin_TransferBetweenAccountTypes_Call{Call: _e.mock.On("TransferBetweenAccountTypes", ctx, fromAccountType, ToAccountType, coin, amount)}
}

func (_c *MockAdmin_TransferBetweenAccountTypes_Call) Run(run func(ctx context.Context, fromAccountType string, ToAccountType string, coin string, amount decimal.Decimal)) *MockAdmin_TransferBetweenAccountTypes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string), args[4].(decimal.Decimal))
	})
	return _c
}

func (_c *MockAdmin_TransferBetweenAccountTypes_Call) Return(_a0 error) *MockAdmin_TransferBetweenAccountTypes_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAdmin_TransferBetweenAccountTypes_Call) RunAndReturn(run func(context.Context, string, string, string, decimal.Decimal) error) *MockAdmin_TransferBetweenAccountTypes_Call {
	_c.Call.Return(run)
	return _c
}

// Withdraw provides a mock function with given fields: ctx, swapID
func (_m *MockAdmin) Withdraw(ctx context.Context, swapID uuid.UUID) (int, error) {
	ret := _m.Called(ctx, swapID)

	if len(ret) == 0 {
		panic("no return value specified for Withdraw")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (int, error)); ok {
		return rf(ctx, swapID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) int); ok {
		r0 = rf(ctx, swapID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, swapID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAdmin_Withdraw_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Withdraw'
type MockAdmin_Withdraw_Call struct {
	*mock.Call
}

// Withdraw is a helper method to define mock.On call
//   - ctx context.Context
//   - swapID uuid.UUID
func (_e *MockAdmin_Expecter) Withdraw(ctx interface{}, swapID interface{}) *MockAdmin_Withdraw_Call {
	return &MockAdmin_Withdraw_Call{Call: _e.mock.On("Withdraw", ctx, swapID)}
}

func (_c *MockAdmin_Withdraw_Call) Run(run func(ctx context.Context, swapID uuid.UUID)) *MockAdmin_Withdraw_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockAdmin_Withdraw_Call) Return(_a0 int, _a1 error) *MockAdmin_Withdraw_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAdmin_Withdraw_Call) RunAndReturn(run func(context.Context, uuid.UUID) (int, error)) *MockAdmin_Withdraw_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAdmin creates a new instance of MockAdmin. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAdmin(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAdmin {
	mock := &MockAdmin{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

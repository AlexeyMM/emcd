// Code generated by mockery v2.43.2. DO NOT EDIT.

package service

import (
	context "context"

	model "code.emcdtech.com/emcd/service/accounting/model"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MockTransaction is an autogenerated mock type for the Transaction type
type MockTransaction struct {
	mock.Mock
}

type MockTransaction_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTransaction) EXPECT() *MockTransaction_Expecter {
	return &MockTransaction_Expecter{mock: &_m.Mock}
}

// ListTransactions provides a mock function with given fields: ctx, receiverAccountIds, types, from, to, limit, fromTransactionID
func (_m *MockTransaction) ListTransactions(ctx context.Context, receiverAccountIds []int, types []int, from time.Time, to time.Time, limit int, fromTransactionID int64) ([]*model.Transaction, int64, error) {
	ret := _m.Called(ctx, receiverAccountIds, types, from, to, limit, fromTransactionID)

	if len(ret) == 0 {
		panic("no return value specified for ListTransactions")
	}

	var r0 []*model.Transaction
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, []int, []int, time.Time, time.Time, int, int64) ([]*model.Transaction, int64, error)); ok {
		return rf(ctx, receiverAccountIds, types, from, to, limit, fromTransactionID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []int, []int, time.Time, time.Time, int, int64) []*model.Transaction); ok {
		r0 = rf(ctx, receiverAccountIds, types, from, to, limit, fromTransactionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []int, []int, time.Time, time.Time, int, int64) int64); ok {
		r1 = rf(ctx, receiverAccountIds, types, from, to, limit, fromTransactionID)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, []int, []int, time.Time, time.Time, int, int64) error); ok {
		r2 = rf(ctx, receiverAccountIds, types, from, to, limit, fromTransactionID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockTransaction_ListTransactions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListTransactions'
type MockTransaction_ListTransactions_Call struct {
	*mock.Call
}

// ListTransactions is a helper method to define mock.On call
//   - ctx context.Context
//   - receiverAccountIds []int
//   - types []int
//   - from time.Time
//   - to time.Time
//   - limit int
//   - fromTransactionID int64
func (_e *MockTransaction_Expecter) ListTransactions(ctx interface{}, receiverAccountIds interface{}, types interface{}, from interface{}, to interface{}, limit interface{}, fromTransactionID interface{}) *MockTransaction_ListTransactions_Call {
	return &MockTransaction_ListTransactions_Call{Call: _e.mock.On("ListTransactions", ctx, receiverAccountIds, types, from, to, limit, fromTransactionID)}
}

func (_c *MockTransaction_ListTransactions_Call) Run(run func(ctx context.Context, receiverAccountIds []int, types []int, from time.Time, to time.Time, limit int, fromTransactionID int64)) *MockTransaction_ListTransactions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]int), args[2].([]int), args[3].(time.Time), args[4].(time.Time), args[5].(int), args[6].(int64))
	})
	return _c
}

func (_c *MockTransaction_ListTransactions_Call) Return(_a0 []*model.Transaction, _a1 int64, _a2 error) *MockTransaction_ListTransactions_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockTransaction_ListTransactions_Call) RunAndReturn(run func(context.Context, []int, []int, time.Time, time.Time, int, int64) ([]*model.Transaction, int64, error)) *MockTransaction_ListTransactions_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTransaction creates a new instance of MockTransaction. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTransaction(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTransaction {
	mock := &MockTransaction{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

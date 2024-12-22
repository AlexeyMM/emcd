// Code generated by mockery v2.46.3. DO NOT EDIT.

package repository

import (
	context "context"

	model "code.emcdtech.com/b2b/swap/model"
	mock "github.com/stretchr/testify/mock"

	pgx "github.com/jackc/pgx/v5"

	transactor "code.emcdtech.com/emcd/sdk/pg"
)

// MockDeposit is an autogenerated mock type for the Deposit type
type MockDeposit struct {
	mock.Mock
}

type MockDeposit_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDeposit) EXPECT() *MockDeposit_Expecter {
	return &MockDeposit_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: ctx, deposit
func (_m *MockDeposit) Add(ctx context.Context, deposit *model.Deposit) error {
	ret := _m.Called(ctx, deposit)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Deposit) error); ok {
		r0 = rf(ctx, deposit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDeposit_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type MockDeposit_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - ctx context.Context
//   - deposit *model.Deposit
func (_e *MockDeposit_Expecter) Add(ctx interface{}, deposit interface{}) *MockDeposit_Add_Call {
	return &MockDeposit_Add_Call{Call: _e.mock.On("Add", ctx, deposit)}
}

func (_c *MockDeposit_Add_Call) Run(run func(ctx context.Context, deposit *model.Deposit)) *MockDeposit_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Deposit))
	})
	return _c
}

func (_c *MockDeposit_Add_Call) Return(_a0 error) *MockDeposit_Add_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDeposit_Add_Call) RunAndReturn(run func(context.Context, *model.Deposit) error) *MockDeposit_Add_Call {
	_c.Call.Return(run)
	return _c
}

// Find provides a mock function with given fields: ctx, filter
func (_m *MockDeposit) Find(ctx context.Context, filter *model.DepositFilter) (model.Deposits, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 model.Deposits
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.DepositFilter) (model.Deposits, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.DepositFilter) model.Deposits); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.Deposits)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.DepositFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDeposit_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type MockDeposit_Find_Call struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - ctx context.Context
//   - filter *model.DepositFilter
func (_e *MockDeposit_Expecter) Find(ctx interface{}, filter interface{}) *MockDeposit_Find_Call {
	return &MockDeposit_Find_Call{Call: _e.mock.On("Find", ctx, filter)}
}

func (_c *MockDeposit_Find_Call) Run(run func(ctx context.Context, filter *model.DepositFilter)) *MockDeposit_Find_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.DepositFilter))
	})
	return _c
}

func (_c *MockDeposit_Find_Call) Return(_a0 model.Deposits, _a1 error) *MockDeposit_Find_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDeposit_Find_Call) RunAndReturn(run func(context.Context, *model.DepositFilter) (model.Deposits, error)) *MockDeposit_Find_Call {
	_c.Call.Return(run)
	return _c
}

// FindOne provides a mock function with given fields: ctx, filter
func (_m *MockDeposit) FindOne(ctx context.Context, filter *model.DepositFilter) (*model.Deposit, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for FindOne")
	}

	var r0 *model.Deposit
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.DepositFilter) (*model.Deposit, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.DepositFilter) *model.Deposit); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Deposit)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.DepositFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDeposit_FindOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindOne'
type MockDeposit_FindOne_Call struct {
	*mock.Call
}

// FindOne is a helper method to define mock.On call
//   - ctx context.Context
//   - filter *model.DepositFilter
func (_e *MockDeposit_Expecter) FindOne(ctx interface{}, filter interface{}) *MockDeposit_FindOne_Call {
	return &MockDeposit_FindOne_Call{Call: _e.mock.On("FindOne", ctx, filter)}
}

func (_c *MockDeposit_FindOne_Call) Run(run func(ctx context.Context, filter *model.DepositFilter)) *MockDeposit_FindOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.DepositFilter))
	})
	return _c
}

func (_c *MockDeposit_FindOne_Call) Return(_a0 *model.Deposit, _a1 error) *MockDeposit_FindOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDeposit_FindOne_Call) RunAndReturn(run func(context.Context, *model.DepositFilter) (*model.Deposit, error)) *MockDeposit_FindOne_Call {
	_c.Call.Return(run)
	return _c
}

// Runner provides a mock function with given fields: ctx
func (_m *MockDeposit) Runner(ctx context.Context) transactor.PgxQueryRunner {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Runner")
	}

	var r0 transactor.PgxQueryRunner
	if rf, ok := ret.Get(0).(func(context.Context) transactor.PgxQueryRunner); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(transactor.PgxQueryRunner)
		}
	}

	return r0
}

// MockDeposit_Runner_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Runner'
type MockDeposit_Runner_Call struct {
	*mock.Call
}

// Runner is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockDeposit_Expecter) Runner(ctx interface{}) *MockDeposit_Runner_Call {
	return &MockDeposit_Runner_Call{Call: _e.mock.On("Runner", ctx)}
}

func (_c *MockDeposit_Runner_Call) Run(run func(ctx context.Context)) *MockDeposit_Runner_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockDeposit_Runner_Call) Return(_a0 transactor.PgxQueryRunner) *MockDeposit_Runner_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDeposit_Runner_Call) RunAndReturn(run func(context.Context) transactor.PgxQueryRunner) *MockDeposit_Runner_Call {
	_c.Call.Return(run)
	return _c
}

// WithinTransaction provides a mock function with given fields: ctx, txFn
func (_m *MockDeposit) WithinTransaction(ctx context.Context, txFn func(context.Context) error) error {
	ret := _m.Called(ctx, txFn)

	if len(ret) == 0 {
		panic("no return value specified for WithinTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context) error) error); ok {
		r0 = rf(ctx, txFn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDeposit_WithinTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithinTransaction'
type MockDeposit_WithinTransaction_Call struct {
	*mock.Call
}

// WithinTransaction is a helper method to define mock.On call
//   - ctx context.Context
//   - txFn func(context.Context) error
func (_e *MockDeposit_Expecter) WithinTransaction(ctx interface{}, txFn interface{}) *MockDeposit_WithinTransaction_Call {
	return &MockDeposit_WithinTransaction_Call{Call: _e.mock.On("WithinTransaction", ctx, txFn)}
}

func (_c *MockDeposit_WithinTransaction_Call) Run(run func(ctx context.Context, txFn func(context.Context) error)) *MockDeposit_WithinTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(context.Context) error))
	})
	return _c
}

func (_c *MockDeposit_WithinTransaction_Call) Return(_a0 error) *MockDeposit_WithinTransaction_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDeposit_WithinTransaction_Call) RunAndReturn(run func(context.Context, func(context.Context) error) error) *MockDeposit_WithinTransaction_Call {
	_c.Call.Return(run)
	return _c
}

// WithinTransactionWithOptions provides a mock function with given fields: ctx, txFn, opts
func (_m *MockDeposit) WithinTransactionWithOptions(ctx context.Context, txFn func(context.Context) error, opts pgx.TxOptions) error {
	ret := _m.Called(ctx, txFn, opts)

	if len(ret) == 0 {
		panic("no return value specified for WithinTransactionWithOptions")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context) error, pgx.TxOptions) error); ok {
		r0 = rf(ctx, txFn, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDeposit_WithinTransactionWithOptions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithinTransactionWithOptions'
type MockDeposit_WithinTransactionWithOptions_Call struct {
	*mock.Call
}

// WithinTransactionWithOptions is a helper method to define mock.On call
//   - ctx context.Context
//   - txFn func(context.Context) error
//   - opts pgx.TxOptions
func (_e *MockDeposit_Expecter) WithinTransactionWithOptions(ctx interface{}, txFn interface{}, opts interface{}) *MockDeposit_WithinTransactionWithOptions_Call {
	return &MockDeposit_WithinTransactionWithOptions_Call{Call: _e.mock.On("WithinTransactionWithOptions", ctx, txFn, opts)}
}

func (_c *MockDeposit_WithinTransactionWithOptions_Call) Run(run func(ctx context.Context, txFn func(context.Context) error, opts pgx.TxOptions)) *MockDeposit_WithinTransactionWithOptions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(context.Context) error), args[2].(pgx.TxOptions))
	})
	return _c
}

func (_c *MockDeposit_WithinTransactionWithOptions_Call) Return(_a0 error) *MockDeposit_WithinTransactionWithOptions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDeposit_WithinTransactionWithOptions_Call) RunAndReturn(run func(context.Context, func(context.Context) error, pgx.TxOptions) error) *MockDeposit_WithinTransactionWithOptions_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDeposit creates a new instance of MockDeposit. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDeposit(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDeposit {
	mock := &MockDeposit{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

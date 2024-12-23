// Code generated by mockery v2.43.2. DO NOT EDIT.

package repository

import (
	context "context"

	model "code.emcdtech.com/b2b/swap/model"
	mock "github.com/stretchr/testify/mock"

	pgx "github.com/jackc/pgx/v5"

	transactor "code.emcdtech.com/emcd/sdk/pg"
)

// MockAccount is an autogenerated mock type for the Account type
type MockAccount struct {
	mock.Mock
}

type MockAccount_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAccount) EXPECT() *MockAccount_Expecter {
	return &MockAccount_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: ctx, account
func (_m *MockAccount) Add(ctx context.Context, account *model.Account) error {
	ret := _m.Called(ctx, account)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Account) error); ok {
		r0 = rf(ctx, account)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAccount_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type MockAccount_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - ctx context.Context
//   - account *model.Account
func (_e *MockAccount_Expecter) Add(ctx interface{}, account interface{}) *MockAccount_Add_Call {
	return &MockAccount_Add_Call{Call: _e.mock.On("Add", ctx, account)}
}

func (_c *MockAccount_Add_Call) Run(run func(ctx context.Context, account *model.Account)) *MockAccount_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Account))
	})
	return _c
}

func (_c *MockAccount_Add_Call) Return(_a0 error) *MockAccount_Add_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAccount_Add_Call) RunAndReturn(run func(context.Context, *model.Account) error) *MockAccount_Add_Call {
	_c.Call.Return(run)
	return _c
}

// Find provides a mock function with given fields: ctx, filter
func (_m *MockAccount) Find(ctx context.Context, filter *model.AccountFilter) (model.Accounts, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 model.Accounts
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AccountFilter) (model.Accounts, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.AccountFilter) model.Accounts); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.Accounts)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.AccountFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAccount_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type MockAccount_Find_Call struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - ctx context.Context
//   - filter *model.AccountFilter
func (_e *MockAccount_Expecter) Find(ctx interface{}, filter interface{}) *MockAccount_Find_Call {
	return &MockAccount_Find_Call{Call: _e.mock.On("Find", ctx, filter)}
}

func (_c *MockAccount_Find_Call) Run(run func(ctx context.Context, filter *model.AccountFilter)) *MockAccount_Find_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.AccountFilter))
	})
	return _c
}

func (_c *MockAccount_Find_Call) Return(_a0 model.Accounts, _a1 error) *MockAccount_Find_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAccount_Find_Call) RunAndReturn(run func(context.Context, *model.AccountFilter) (model.Accounts, error)) *MockAccount_Find_Call {
	_c.Call.Return(run)
	return _c
}

// FindOne provides a mock function with given fields: ctx, filter
func (_m *MockAccount) FindOne(ctx context.Context, filter *model.AccountFilter) (*model.Account, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for FindOne")
	}

	var r0 *model.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AccountFilter) (*model.Account, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.AccountFilter) *model.Account); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.AccountFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAccount_FindOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindOne'
type MockAccount_FindOne_Call struct {
	*mock.Call
}

// FindOne is a helper method to define mock.On call
//   - ctx context.Context
//   - filter *model.AccountFilter
func (_e *MockAccount_Expecter) FindOne(ctx interface{}, filter interface{}) *MockAccount_FindOne_Call {
	return &MockAccount_FindOne_Call{Call: _e.mock.On("FindOne", ctx, filter)}
}

func (_c *MockAccount_FindOne_Call) Run(run func(ctx context.Context, filter *model.AccountFilter)) *MockAccount_FindOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.AccountFilter))
	})
	return _c
}

func (_c *MockAccount_FindOne_Call) Return(_a0 *model.Account, _a1 error) *MockAccount_FindOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAccount_FindOne_Call) RunAndReturn(run func(context.Context, *model.AccountFilter) (*model.Account, error)) *MockAccount_FindOne_Call {
	_c.Call.Return(run)
	return _c
}

// Runner provides a mock function with given fields: ctx
func (_m *MockAccount) Runner(ctx context.Context) transactor.PgxQueryRunner {
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

// MockAccount_Runner_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Runner'
type MockAccount_Runner_Call struct {
	*mock.Call
}

// Runner is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockAccount_Expecter) Runner(ctx interface{}) *MockAccount_Runner_Call {
	return &MockAccount_Runner_Call{Call: _e.mock.On("Runner", ctx)}
}

func (_c *MockAccount_Runner_Call) Run(run func(ctx context.Context)) *MockAccount_Runner_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockAccount_Runner_Call) Return(_a0 transactor.PgxQueryRunner) *MockAccount_Runner_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAccount_Runner_Call) RunAndReturn(run func(context.Context) transactor.PgxQueryRunner) *MockAccount_Runner_Call {
	_c.Call.Return(run)
	return _c
}

// WithinTransaction provides a mock function with given fields: ctx, txFn
func (_m *MockAccount) WithinTransaction(ctx context.Context, txFn func(context.Context) error) error {
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

// MockAccount_WithinTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithinTransaction'
type MockAccount_WithinTransaction_Call struct {
	*mock.Call
}

// WithinTransaction is a helper method to define mock.On call
//   - ctx context.Context
//   - txFn func(context.Context) error
func (_e *MockAccount_Expecter) WithinTransaction(ctx interface{}, txFn interface{}) *MockAccount_WithinTransaction_Call {
	return &MockAccount_WithinTransaction_Call{Call: _e.mock.On("WithinTransaction", ctx, txFn)}
}

func (_c *MockAccount_WithinTransaction_Call) Run(run func(ctx context.Context, txFn func(context.Context) error)) *MockAccount_WithinTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(context.Context) error))
	})
	return _c
}

func (_c *MockAccount_WithinTransaction_Call) Return(_a0 error) *MockAccount_WithinTransaction_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAccount_WithinTransaction_Call) RunAndReturn(run func(context.Context, func(context.Context) error) error) *MockAccount_WithinTransaction_Call {
	_c.Call.Return(run)
	return _c
}

// WithinTransactionWithOptions provides a mock function with given fields: ctx, txFn, opts
func (_m *MockAccount) WithinTransactionWithOptions(ctx context.Context, txFn func(context.Context) error, opts pgx.TxOptions) error {
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

// MockAccount_WithinTransactionWithOptions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithinTransactionWithOptions'
type MockAccount_WithinTransactionWithOptions_Call struct {
	*mock.Call
}

// WithinTransactionWithOptions is a helper method to define mock.On call
//   - ctx context.Context
//   - txFn func(context.Context) error
//   - opts pgx.TxOptions
func (_e *MockAccount_Expecter) WithinTransactionWithOptions(ctx interface{}, txFn interface{}, opts interface{}) *MockAccount_WithinTransactionWithOptions_Call {
	return &MockAccount_WithinTransactionWithOptions_Call{Call: _e.mock.On("WithinTransactionWithOptions", ctx, txFn, opts)}
}

func (_c *MockAccount_WithinTransactionWithOptions_Call) Run(run func(ctx context.Context, txFn func(context.Context) error, opts pgx.TxOptions)) *MockAccount_WithinTransactionWithOptions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(context.Context) error), args[2].(pgx.TxOptions))
	})
	return _c
}

func (_c *MockAccount_WithinTransactionWithOptions_Call) Return(_a0 error) *MockAccount_WithinTransactionWithOptions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAccount_WithinTransactionWithOptions_Call) RunAndReturn(run func(context.Context, func(context.Context) error, pgx.TxOptions) error) *MockAccount_WithinTransactionWithOptions_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAccount creates a new instance of MockAccount. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAccount(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAccount {
	mock := &MockAccount{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Code generated by mockery v2.46.3. DO NOT EDIT.

package repository

import (
	context "context"

	model "code.emcdtech.com/b2b/swap/model"
	mock "github.com/stretchr/testify/mock"

	pgx "github.com/jackc/pgx/v5"

	transactor "code.emcdtech.com/emcd/sdk/pg"
)

// MockSwap is an autogenerated mock type for the Swap type
type MockSwap struct {
	mock.Mock
}

type MockSwap_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSwap) EXPECT() *MockSwap_Expecter {
	return &MockSwap_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: ctx, swap
func (_m *MockSwap) Add(ctx context.Context, swap *model.Swap) error {
	ret := _m.Called(ctx, swap)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Swap) error); ok {
		r0 = rf(ctx, swap)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSwap_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type MockSwap_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - ctx context.Context
//   - swap *model.Swap
func (_e *MockSwap_Expecter) Add(ctx interface{}, swap interface{}) *MockSwap_Add_Call {
	return &MockSwap_Add_Call{Call: _e.mock.On("Add", ctx, swap)}
}

func (_c *MockSwap_Add_Call) Run(run func(ctx context.Context, swap *model.Swap)) *MockSwap_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Swap))
	})
	return _c
}

func (_c *MockSwap_Add_Call) Return(_a0 error) *MockSwap_Add_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSwap_Add_Call) RunAndReturn(run func(context.Context, *model.Swap) error) *MockSwap_Add_Call {
	_c.Call.Return(run)
	return _c
}

// CountSwapsByStatus provides a mock function with given fields: ctx
func (_m *MockSwap) CountSwapsByStatus(ctx context.Context) (map[model.Status]int, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for CountSwapsByStatus")
	}

	var r0 map[model.Status]int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (map[model.Status]int, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) map[model.Status]int); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[model.Status]int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSwap_CountSwapsByStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CountSwapsByStatus'
type MockSwap_CountSwapsByStatus_Call struct {
	*mock.Call
}

// CountSwapsByStatus is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockSwap_Expecter) CountSwapsByStatus(ctx interface{}) *MockSwap_CountSwapsByStatus_Call {
	return &MockSwap_CountSwapsByStatus_Call{Call: _e.mock.On("CountSwapsByStatus", ctx)}
}

func (_c *MockSwap_CountSwapsByStatus_Call) Run(run func(ctx context.Context)) *MockSwap_CountSwapsByStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockSwap_CountSwapsByStatus_Call) Return(_a0 map[model.Status]int, _a1 error) *MockSwap_CountSwapsByStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSwap_CountSwapsByStatus_Call) RunAndReturn(run func(context.Context) (map[model.Status]int, error)) *MockSwap_CountSwapsByStatus_Call {
	_c.Call.Return(run)
	return _c
}

// CountTotalWithFilter provides a mock function with given fields: ctx, filter
func (_m *MockSwap) CountTotalWithFilter(ctx context.Context, filter *model.SwapFilter) (int, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for CountTotalWithFilter")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.SwapFilter) (int, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.SwapFilter) int); ok {
		r0 = rf(ctx, filter)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.SwapFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSwap_CountTotalWithFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CountTotalWithFilter'
type MockSwap_CountTotalWithFilter_Call struct {
	*mock.Call
}

// CountTotalWithFilter is a helper method to define mock.On call
//   - ctx context.Context
//   - filter *model.SwapFilter
func (_e *MockSwap_Expecter) CountTotalWithFilter(ctx interface{}, filter interface{}) *MockSwap_CountTotalWithFilter_Call {
	return &MockSwap_CountTotalWithFilter_Call{Call: _e.mock.On("CountTotalWithFilter", ctx, filter)}
}

func (_c *MockSwap_CountTotalWithFilter_Call) Run(run func(ctx context.Context, filter *model.SwapFilter)) *MockSwap_CountTotalWithFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.SwapFilter))
	})
	return _c
}

func (_c *MockSwap_CountTotalWithFilter_Call) Return(_a0 int, _a1 error) *MockSwap_CountTotalWithFilter_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSwap_CountTotalWithFilter_Call) RunAndReturn(run func(context.Context, *model.SwapFilter) (int, error)) *MockSwap_CountTotalWithFilter_Call {
	_c.Call.Return(run)
	return _c
}

// Find provides a mock function with given fields: ctx, filter
func (_m *MockSwap) Find(ctx context.Context, filter *model.SwapFilter) (model.Swaps, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 model.Swaps
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.SwapFilter) (model.Swaps, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.SwapFilter) model.Swaps); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.Swaps)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.SwapFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSwap_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type MockSwap_Find_Call struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - ctx context.Context
//   - filter *model.SwapFilter
func (_e *MockSwap_Expecter) Find(ctx interface{}, filter interface{}) *MockSwap_Find_Call {
	return &MockSwap_Find_Call{Call: _e.mock.On("Find", ctx, filter)}
}

func (_c *MockSwap_Find_Call) Run(run func(ctx context.Context, filter *model.SwapFilter)) *MockSwap_Find_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.SwapFilter))
	})
	return _c
}

func (_c *MockSwap_Find_Call) Return(_a0 model.Swaps, _a1 error) *MockSwap_Find_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSwap_Find_Call) RunAndReturn(run func(context.Context, *model.SwapFilter) (model.Swaps, error)) *MockSwap_Find_Call {
	_c.Call.Return(run)
	return _c
}

// FindOne provides a mock function with given fields: ctx, filter
func (_m *MockSwap) FindOne(ctx context.Context, filter *model.SwapFilter) (*model.Swap, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for FindOne")
	}

	var r0 *model.Swap
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.SwapFilter) (*model.Swap, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.SwapFilter) *model.Swap); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Swap)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.SwapFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSwap_FindOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindOne'
type MockSwap_FindOne_Call struct {
	*mock.Call
}

// FindOne is a helper method to define mock.On call
//   - ctx context.Context
//   - filter *model.SwapFilter
func (_e *MockSwap_Expecter) FindOne(ctx interface{}, filter interface{}) *MockSwap_FindOne_Call {
	return &MockSwap_FindOne_Call{Call: _e.mock.On("FindOne", ctx, filter)}
}

func (_c *MockSwap_FindOne_Call) Run(run func(ctx context.Context, filter *model.SwapFilter)) *MockSwap_FindOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.SwapFilter))
	})
	return _c
}

func (_c *MockSwap_FindOne_Call) Return(_a0 *model.Swap, _a1 error) *MockSwap_FindOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSwap_FindOne_Call) RunAndReturn(run func(context.Context, *model.SwapFilter) (*model.Swap, error)) *MockSwap_FindOne_Call {
	_c.Call.Return(run)
	return _c
}

// Runner provides a mock function with given fields: ctx
func (_m *MockSwap) Runner(ctx context.Context) transactor.PgxQueryRunner {
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

// MockSwap_Runner_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Runner'
type MockSwap_Runner_Call struct {
	*mock.Call
}

// Runner is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockSwap_Expecter) Runner(ctx interface{}) *MockSwap_Runner_Call {
	return &MockSwap_Runner_Call{Call: _e.mock.On("Runner", ctx)}
}

func (_c *MockSwap_Runner_Call) Run(run func(ctx context.Context)) *MockSwap_Runner_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockSwap_Runner_Call) Return(_a0 transactor.PgxQueryRunner) *MockSwap_Runner_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSwap_Runner_Call) RunAndReturn(run func(context.Context) transactor.PgxQueryRunner) *MockSwap_Runner_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, swap, filter, partial
func (_m *MockSwap) Update(ctx context.Context, swap *model.Swap, filter *model.SwapFilter, partial *model.SwapPartial) error {
	ret := _m.Called(ctx, swap, filter, partial)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Swap, *model.SwapFilter, *model.SwapPartial) error); ok {
		r0 = rf(ctx, swap, filter, partial)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSwap_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockSwap_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - swap *model.Swap
//   - filter *model.SwapFilter
//   - partial *model.SwapPartial
func (_e *MockSwap_Expecter) Update(ctx interface{}, swap interface{}, filter interface{}, partial interface{}) *MockSwap_Update_Call {
	return &MockSwap_Update_Call{Call: _e.mock.On("Update", ctx, swap, filter, partial)}
}

func (_c *MockSwap_Update_Call) Run(run func(ctx context.Context, swap *model.Swap, filter *model.SwapFilter, partial *model.SwapPartial)) *MockSwap_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Swap), args[2].(*model.SwapFilter), args[3].(*model.SwapPartial))
	})
	return _c
}

func (_c *MockSwap_Update_Call) Return(_a0 error) *MockSwap_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSwap_Update_Call) RunAndReturn(run func(context.Context, *model.Swap, *model.SwapFilter, *model.SwapPartial) error) *MockSwap_Update_Call {
	_c.Call.Return(run)
	return _c
}

// WithinTransaction provides a mock function with given fields: ctx, txFn
func (_m *MockSwap) WithinTransaction(ctx context.Context, txFn func(context.Context) error) error {
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

// MockSwap_WithinTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithinTransaction'
type MockSwap_WithinTransaction_Call struct {
	*mock.Call
}

// WithinTransaction is a helper method to define mock.On call
//   - ctx context.Context
//   - txFn func(context.Context) error
func (_e *MockSwap_Expecter) WithinTransaction(ctx interface{}, txFn interface{}) *MockSwap_WithinTransaction_Call {
	return &MockSwap_WithinTransaction_Call{Call: _e.mock.On("WithinTransaction", ctx, txFn)}
}

func (_c *MockSwap_WithinTransaction_Call) Run(run func(ctx context.Context, txFn func(context.Context) error)) *MockSwap_WithinTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(context.Context) error))
	})
	return _c
}

func (_c *MockSwap_WithinTransaction_Call) Return(_a0 error) *MockSwap_WithinTransaction_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSwap_WithinTransaction_Call) RunAndReturn(run func(context.Context, func(context.Context) error) error) *MockSwap_WithinTransaction_Call {
	_c.Call.Return(run)
	return _c
}

// WithinTransactionWithOptions provides a mock function with given fields: ctx, txFn, opts
func (_m *MockSwap) WithinTransactionWithOptions(ctx context.Context, txFn func(context.Context) error, opts pgx.TxOptions) error {
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

// MockSwap_WithinTransactionWithOptions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithinTransactionWithOptions'
type MockSwap_WithinTransactionWithOptions_Call struct {
	*mock.Call
}

// WithinTransactionWithOptions is a helper method to define mock.On call
//   - ctx context.Context
//   - txFn func(context.Context) error
//   - opts pgx.TxOptions
func (_e *MockSwap_Expecter) WithinTransactionWithOptions(ctx interface{}, txFn interface{}, opts interface{}) *MockSwap_WithinTransactionWithOptions_Call {
	return &MockSwap_WithinTransactionWithOptions_Call{Call: _e.mock.On("WithinTransactionWithOptions", ctx, txFn, opts)}
}

func (_c *MockSwap_WithinTransactionWithOptions_Call) Run(run func(ctx context.Context, txFn func(context.Context) error, opts pgx.TxOptions)) *MockSwap_WithinTransactionWithOptions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(context.Context) error), args[2].(pgx.TxOptions))
	})
	return _c
}

func (_c *MockSwap_WithinTransactionWithOptions_Call) Return(_a0 error) *MockSwap_WithinTransactionWithOptions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSwap_WithinTransactionWithOptions_Call) RunAndReturn(run func(context.Context, func(context.Context) error, pgx.TxOptions) error) *MockSwap_WithinTransactionWithOptions_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSwap creates a new instance of MockSwap. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSwap(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSwap {
	mock := &MockSwap{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
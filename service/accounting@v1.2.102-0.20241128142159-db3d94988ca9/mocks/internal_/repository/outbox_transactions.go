// Code generated by mockery v2.43.2. DO NOT EDIT.

package repository

import (
	context "context"

	pgx "github.com/jackc/pgx/v5"
	mock "github.com/stretchr/testify/mock"
)

// MockOutboxTransactions is an autogenerated mock type for the OutboxTransactions type
type MockOutboxTransactions struct {
	mock.Mock
}

type MockOutboxTransactions_Expecter struct {
	mock *mock.Mock
}

func (_m *MockOutboxTransactions) EXPECT() *MockOutboxTransactions_Expecter {
	return &MockOutboxTransactions_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: ctx, tx, ids
func (_m *MockOutboxTransactions) Delete(ctx context.Context, tx pgx.Tx, ids ...int64) error {
	_va := make([]interface{}, len(ids))
	for _i := range ids {
		_va[_i] = ids[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, tx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx, ...int64) error); ok {
		r0 = rf(ctx, tx, ids...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockOutboxTransactions_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockOutboxTransactions_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - tx pgx.Tx
//   - ids ...int64
func (_e *MockOutboxTransactions_Expecter) Delete(ctx interface{}, tx interface{}, ids ...interface{}) *MockOutboxTransactions_Delete_Call {
	return &MockOutboxTransactions_Delete_Call{Call: _e.mock.On("Delete",
		append([]interface{}{ctx, tx}, ids...)...)}
}

func (_c *MockOutboxTransactions_Delete_Call) Run(run func(ctx context.Context, tx pgx.Tx, ids ...int64)) *MockOutboxTransactions_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]int64, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(int64)
			}
		}
		run(args[0].(context.Context), args[1].(pgx.Tx), variadicArgs...)
	})
	return _c
}

func (_c *MockOutboxTransactions_Delete_Call) Return(_a0 error) *MockOutboxTransactions_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockOutboxTransactions_Delete_Call) RunAndReturn(run func(context.Context, pgx.Tx, ...int64) error) *MockOutboxTransactions_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, tx, limit
func (_m *MockOutboxTransactions) List(ctx context.Context, tx pgx.Tx, limit uint) ([]int64, error) {
	ret := _m.Called(ctx, tx, limit)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx, uint) ([]int64, error)); ok {
		return rf(ctx, tx, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx, uint) []int64); ok {
		r0 = rf(ctx, tx, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int64)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, pgx.Tx, uint) error); ok {
		r1 = rf(ctx, tx, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOutboxTransactions_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type MockOutboxTransactions_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - tx pgx.Tx
//   - limit uint
func (_e *MockOutboxTransactions_Expecter) List(ctx interface{}, tx interface{}, limit interface{}) *MockOutboxTransactions_List_Call {
	return &MockOutboxTransactions_List_Call{Call: _e.mock.On("List", ctx, tx, limit)}
}

func (_c *MockOutboxTransactions_List_Call) Run(run func(ctx context.Context, tx pgx.Tx, limit uint)) *MockOutboxTransactions_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(pgx.Tx), args[2].(uint))
	})
	return _c
}

func (_c *MockOutboxTransactions_List_Call) Return(_a0 []int64, _a1 error) *MockOutboxTransactions_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOutboxTransactions_List_Call) RunAndReturn(run func(context.Context, pgx.Tx, uint) ([]int64, error)) *MockOutboxTransactions_List_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockOutboxTransactions creates a new instance of MockOutboxTransactions. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockOutboxTransactions(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockOutboxTransactions {
	mock := &MockOutboxTransactions{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

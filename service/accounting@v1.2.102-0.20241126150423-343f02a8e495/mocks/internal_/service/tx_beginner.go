// Code generated by mockery v2.43.2. DO NOT EDIT.

package service

import (
	context "context"

	pgx "github.com/jackc/pgx/v5"
	mock "github.com/stretchr/testify/mock"
)

// MockTxBeginner is an autogenerated mock type for the TxBeginner type
type MockTxBeginner struct {
	mock.Mock
}

type MockTxBeginner_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTxBeginner) EXPECT() *MockTxBeginner_Expecter {
	return &MockTxBeginner_Expecter{mock: &_m.Mock}
}

// Begin provides a mock function with given fields: ctx
func (_m *MockTxBeginner) Begin(ctx context.Context) (pgx.Tx, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Begin")
	}

	var r0 pgx.Tx
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (pgx.Tx, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) pgx.Tx); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgx.Tx)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTxBeginner_Begin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Begin'
type MockTxBeginner_Begin_Call struct {
	*mock.Call
}

// Begin is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockTxBeginner_Expecter) Begin(ctx interface{}) *MockTxBeginner_Begin_Call {
	return &MockTxBeginner_Begin_Call{Call: _e.mock.On("Begin", ctx)}
}

func (_c *MockTxBeginner_Begin_Call) Run(run func(ctx context.Context)) *MockTxBeginner_Begin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockTxBeginner_Begin_Call) Return(_a0 pgx.Tx, _a1 error) *MockTxBeginner_Begin_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTxBeginner_Begin_Call) RunAndReturn(run func(context.Context) (pgx.Tx, error)) *MockTxBeginner_Begin_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTxBeginner creates a new instance of MockTxBeginner. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTxBeginner(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTxBeginner {
	mock := &MockTxBeginner{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

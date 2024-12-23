// Code generated by mockery v2.43.2. DO NOT EDIT.

package repository

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MockLimitPayouts is an autogenerated mock type for the LimitPayouts type
type MockLimitPayouts struct {
	mock.Mock
}

type MockLimitPayouts_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLimitPayouts) EXPECT() *MockLimitPayouts_Expecter {
	return &MockLimitPayouts_Expecter{mock: &_m.Mock}
}

// GetLimit provides a mock function with given fields: ctx, coinID
func (_m *MockLimitPayouts) GetLimit(ctx context.Context, coinID int) (float64, error) {
	ret := _m.Called(ctx, coinID)

	if len(ret) == 0 {
		panic("no return value specified for GetLimit")
	}

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (float64, error)); ok {
		return rf(ctx, coinID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) float64); ok {
		r0 = rf(ctx, coinID)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, coinID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockLimitPayouts_GetLimit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLimit'
type MockLimitPayouts_GetLimit_Call struct {
	*mock.Call
}

// GetLimit is a helper method to define mock.On call
//   - ctx context.Context
//   - coinID int
func (_e *MockLimitPayouts_Expecter) GetLimit(ctx interface{}, coinID interface{}) *MockLimitPayouts_GetLimit_Call {
	return &MockLimitPayouts_GetLimit_Call{Call: _e.mock.On("GetLimit", ctx, coinID)}
}

func (_c *MockLimitPayouts_GetLimit_Call) Run(run func(ctx context.Context, coinID int)) *MockLimitPayouts_GetLimit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *MockLimitPayouts_GetLimit_Call) Return(_a0 float64, _a1 error) *MockLimitPayouts_GetLimit_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockLimitPayouts_GetLimit_Call) RunAndReturn(run func(context.Context, int) (float64, error)) *MockLimitPayouts_GetLimit_Call {
	_c.Call.Return(run)
	return _c
}

// GetMainUserId provides a mock function with given fields: ctx, userID
func (_m *MockLimitPayouts) GetMainUserId(ctx context.Context, userID int64) (int64, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetMainUserId")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (int64, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) int64); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockLimitPayouts_GetMainUserId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMainUserId'
type MockLimitPayouts_GetMainUserId_Call struct {
	*mock.Call
}

// GetMainUserId is a helper method to define mock.On call
//   - ctx context.Context
//   - userID int64
func (_e *MockLimitPayouts_Expecter) GetMainUserId(ctx interface{}, userID interface{}) *MockLimitPayouts_GetMainUserId_Call {
	return &MockLimitPayouts_GetMainUserId_Call{Call: _e.mock.On("GetMainUserId", ctx, userID)}
}

func (_c *MockLimitPayouts_GetMainUserId_Call) Run(run func(ctx context.Context, userID int64)) *MockLimitPayouts_GetMainUserId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockLimitPayouts_GetMainUserId_Call) Return(_a0 int64, _a1 error) *MockLimitPayouts_GetMainUserId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockLimitPayouts_GetMainUserId_Call) RunAndReturn(run func(context.Context, int64) (int64, error)) *MockLimitPayouts_GetMainUserId_Call {
	_c.Call.Return(run)
	return _c
}

// GetRedisBlockStatus provides a mock function with given fields: ctx, userID
func (_m *MockLimitPayouts) GetRedisBlockStatus(ctx context.Context, userID int64) (int, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetRedisBlockStatus")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (int, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) int); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockLimitPayouts_GetRedisBlockStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRedisBlockStatus'
type MockLimitPayouts_GetRedisBlockStatus_Call struct {
	*mock.Call
}

// GetRedisBlockStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - userID int64
func (_e *MockLimitPayouts_Expecter) GetRedisBlockStatus(ctx interface{}, userID interface{}) *MockLimitPayouts_GetRedisBlockStatus_Call {
	return &MockLimitPayouts_GetRedisBlockStatus_Call{Call: _e.mock.On("GetRedisBlockStatus", ctx, userID)}
}

func (_c *MockLimitPayouts_GetRedisBlockStatus_Call) Run(run func(ctx context.Context, userID int64)) *MockLimitPayouts_GetRedisBlockStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockLimitPayouts_GetRedisBlockStatus_Call) Return(_a0 int, _a1 error) *MockLimitPayouts_GetRedisBlockStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockLimitPayouts_GetRedisBlockStatus_Call) RunAndReturn(run func(context.Context, int64) (int, error)) *MockLimitPayouts_GetRedisBlockStatus_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserPayoutsSum provides a mock function with given fields: ctx, userID, coinId
func (_m *MockLimitPayouts) GetUserPayoutsSum(ctx context.Context, userID int64, coinId int) (float64, error) {
	ret := _m.Called(ctx, userID, coinId)

	if len(ret) == 0 {
		panic("no return value specified for GetUserPayoutsSum")
	}

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int) (float64, error)); ok {
		return rf(ctx, userID, coinId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int) float64); ok {
		r0 = rf(ctx, userID, coinId)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int) error); ok {
		r1 = rf(ctx, userID, coinId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockLimitPayouts_GetUserPayoutsSum_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserPayoutsSum'
type MockLimitPayouts_GetUserPayoutsSum_Call struct {
	*mock.Call
}

// GetUserPayoutsSum is a helper method to define mock.On call
//   - ctx context.Context
//   - userID int64
//   - coinId int
func (_e *MockLimitPayouts_Expecter) GetUserPayoutsSum(ctx interface{}, userID interface{}, coinId interface{}) *MockLimitPayouts_GetUserPayoutsSum_Call {
	return &MockLimitPayouts_GetUserPayoutsSum_Call{Call: _e.mock.On("GetUserPayoutsSum", ctx, userID, coinId)}
}

func (_c *MockLimitPayouts_GetUserPayoutsSum_Call) Run(run func(ctx context.Context, userID int64, coinId int)) *MockLimitPayouts_GetUserPayoutsSum_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(int))
	})
	return _c
}

func (_c *MockLimitPayouts_GetUserPayoutsSum_Call) Return(_a0 float64, _a1 error) *MockLimitPayouts_GetUserPayoutsSum_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockLimitPayouts_GetUserPayoutsSum_Call) RunAndReturn(run func(context.Context, int64, int) (float64, error)) *MockLimitPayouts_GetUserPayoutsSum_Call {
	_c.Call.Return(run)
	return _c
}

// SetRedisBlockStatus provides a mock function with given fields: ctx, userID, status, exp
func (_m *MockLimitPayouts) SetRedisBlockStatus(ctx context.Context, userID int64, status int, exp time.Duration) error {
	ret := _m.Called(ctx, userID, status, exp)

	if len(ret) == 0 {
		panic("no return value specified for SetRedisBlockStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int, time.Duration) error); ok {
		r0 = rf(ctx, userID, status, exp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockLimitPayouts_SetRedisBlockStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetRedisBlockStatus'
type MockLimitPayouts_SetRedisBlockStatus_Call struct {
	*mock.Call
}

// SetRedisBlockStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - userID int64
//   - status int
//   - exp time.Duration
func (_e *MockLimitPayouts_Expecter) SetRedisBlockStatus(ctx interface{}, userID interface{}, status interface{}, exp interface{}) *MockLimitPayouts_SetRedisBlockStatus_Call {
	return &MockLimitPayouts_SetRedisBlockStatus_Call{Call: _e.mock.On("SetRedisBlockStatus", ctx, userID, status, exp)}
}

func (_c *MockLimitPayouts_SetRedisBlockStatus_Call) Run(run func(ctx context.Context, userID int64, status int, exp time.Duration)) *MockLimitPayouts_SetRedisBlockStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(int), args[3].(time.Duration))
	})
	return _c
}

func (_c *MockLimitPayouts_SetRedisBlockStatus_Call) Return(_a0 error) *MockLimitPayouts_SetRedisBlockStatus_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLimitPayouts_SetRedisBlockStatus_Call) RunAndReturn(run func(context.Context, int64, int, time.Duration) error) *MockLimitPayouts_SetRedisBlockStatus_Call {
	_c.Call.Return(run)
	return _c
}

// SetUserNopay provides a mock function with given fields: ctx, userID
func (_m *MockLimitPayouts) SetUserNopay(ctx context.Context, userID int64) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for SetUserNopay")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockLimitPayouts_SetUserNopay_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetUserNopay'
type MockLimitPayouts_SetUserNopay_Call struct {
	*mock.Call
}

// SetUserNopay is a helper method to define mock.On call
//   - ctx context.Context
//   - userID int64
func (_e *MockLimitPayouts_Expecter) SetUserNopay(ctx interface{}, userID interface{}) *MockLimitPayouts_SetUserNopay_Call {
	return &MockLimitPayouts_SetUserNopay_Call{Call: _e.mock.On("SetUserNopay", ctx, userID)}
}

func (_c *MockLimitPayouts_SetUserNopay_Call) Run(run func(ctx context.Context, userID int64)) *MockLimitPayouts_SetUserNopay_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockLimitPayouts_SetUserNopay_Call) Return(_a0 error) *MockLimitPayouts_SetUserNopay_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLimitPayouts_SetUserNopay_Call) RunAndReturn(run func(context.Context, int64) error) *MockLimitPayouts_SetUserNopay_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockLimitPayouts creates a new instance of MockLimitPayouts. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockLimitPayouts(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockLimitPayouts {
	mock := &MockLimitPayouts{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	sync "sync"

	time "time"
)

// MockCoinValidatorRepository is an autogenerated mock type for the CoinValidatorRepository type
type MockCoinValidatorRepository struct {
	mock.Mock
}

type MockCoinValidatorRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCoinValidatorRepository) EXPECT() *MockCoinValidatorRepository_Expecter {
	return &MockCoinValidatorRepository_Expecter{mock: &_m.Mock}
}

// GetCodeById provides a mock function with given fields: coinIdLegacy
func (_m *MockCoinValidatorRepository) GetCodeById(coinIdLegacy int32) (string, bool) {
	ret := _m.Called(coinIdLegacy)

	if len(ret) == 0 {
		panic("no return value specified for GetCodeById")
	}

	var r0 string
	var r1 bool
	if rf, ok := ret.Get(0).(func(int32) (string, bool)); ok {
		return rf(coinIdLegacy)
	}
	if rf, ok := ret.Get(0).(func(int32) string); ok {
		r0 = rf(coinIdLegacy)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(int32) bool); ok {
		r1 = rf(coinIdLegacy)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// MockCoinValidatorRepository_GetCodeById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCodeById'
type MockCoinValidatorRepository_GetCodeById_Call struct {
	*mock.Call
}

// GetCodeById is a helper method to define mock.On call
//   - coinIdLegacy int32
func (_e *MockCoinValidatorRepository_Expecter) GetCodeById(coinIdLegacy interface{}) *MockCoinValidatorRepository_GetCodeById_Call {
	return &MockCoinValidatorRepository_GetCodeById_Call{Call: _e.mock.On("GetCodeById", coinIdLegacy)}
}

func (_c *MockCoinValidatorRepository_GetCodeById_Call) Run(run func(coinIdLegacy int32)) *MockCoinValidatorRepository_GetCodeById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int32))
	})
	return _c
}

func (_c *MockCoinValidatorRepository_GetCodeById_Call) Return(_a0 string, _a1 bool) *MockCoinValidatorRepository_GetCodeById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCoinValidatorRepository_GetCodeById_Call) RunAndReturn(run func(int32) (string, bool)) *MockCoinValidatorRepository_GetCodeById_Call {
	_c.Call.Return(run)
	return _c
}

// GetCodes provides a mock function with given fields:
func (_m *MockCoinValidatorRepository) GetCodes() []string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetCodes")
	}

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// MockCoinValidatorRepository_GetCodes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCodes'
type MockCoinValidatorRepository_GetCodes_Call struct {
	*mock.Call
}

// GetCodes is a helper method to define mock.On call
func (_e *MockCoinValidatorRepository_Expecter) GetCodes() *MockCoinValidatorRepository_GetCodes_Call {
	return &MockCoinValidatorRepository_GetCodes_Call{Call: _e.mock.On("GetCodes")}
}

func (_c *MockCoinValidatorRepository_GetCodes_Call) Run(run func()) *MockCoinValidatorRepository_GetCodes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCoinValidatorRepository_GetCodes_Call) Return(_a0 []string) *MockCoinValidatorRepository_GetCodes_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCoinValidatorRepository_GetCodes_Call) RunAndReturn(run func() []string) *MockCoinValidatorRepository_GetCodes_Call {
	_c.Call.Return(run)
	return _c
}

// GetIdByCode provides a mock function with given fields: coinCode
func (_m *MockCoinValidatorRepository) GetIdByCode(coinCode string) (int32, bool) {
	ret := _m.Called(coinCode)

	if len(ret) == 0 {
		panic("no return value specified for GetIdByCode")
	}

	var r0 int32
	var r1 bool
	if rf, ok := ret.Get(0).(func(string) (int32, bool)); ok {
		return rf(coinCode)
	}
	if rf, ok := ret.Get(0).(func(string) int32); ok {
		r0 = rf(coinCode)
	} else {
		r0 = ret.Get(0).(int32)
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(coinCode)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// MockCoinValidatorRepository_GetIdByCode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIdByCode'
type MockCoinValidatorRepository_GetIdByCode_Call struct {
	*mock.Call
}

// GetIdByCode is a helper method to define mock.On call
//   - coinCode string
func (_e *MockCoinValidatorRepository_Expecter) GetIdByCode(coinCode interface{}) *MockCoinValidatorRepository_GetIdByCode_Call {
	return &MockCoinValidatorRepository_GetIdByCode_Call{Call: _e.mock.On("GetIdByCode", coinCode)}
}

func (_c *MockCoinValidatorRepository_GetIdByCode_Call) Run(run func(coinCode string)) *MockCoinValidatorRepository_GetIdByCode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCoinValidatorRepository_GetIdByCode_Call) Return(_a0 int32, _a1 bool) *MockCoinValidatorRepository_GetIdByCode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCoinValidatorRepository_GetIdByCode_Call) RunAndReturn(run func(string) (int32, bool)) *MockCoinValidatorRepository_GetIdByCode_Call {
	_c.Call.Return(run)
	return _c
}

// GetIdsLegacy provides a mock function with given fields:
func (_m *MockCoinValidatorRepository) GetIdsLegacy() []int32 {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetIdsLegacy")
	}

	var r0 []int32
	if rf, ok := ret.Get(0).(func() []int32); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int32)
		}
	}

	return r0
}

// MockCoinValidatorRepository_GetIdsLegacy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIdsLegacy'
type MockCoinValidatorRepository_GetIdsLegacy_Call struct {
	*mock.Call
}

// GetIdsLegacy is a helper method to define mock.On call
func (_e *MockCoinValidatorRepository_Expecter) GetIdsLegacy() *MockCoinValidatorRepository_GetIdsLegacy_Call {
	return &MockCoinValidatorRepository_GetIdsLegacy_Call{Call: _e.mock.On("GetIdsLegacy")}
}

func (_c *MockCoinValidatorRepository_GetIdsLegacy_Call) Run(run func()) *MockCoinValidatorRepository_GetIdsLegacy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCoinValidatorRepository_GetIdsLegacy_Call) Return(_a0 []int32) *MockCoinValidatorRepository_GetIdsLegacy_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCoinValidatorRepository_GetIdsLegacy_Call) RunAndReturn(run func() []int32) *MockCoinValidatorRepository_GetIdsLegacy_Call {
	_c.Call.Return(run)
	return _c
}

// IsValidCode provides a mock function with given fields: coinCode
func (_m *MockCoinValidatorRepository) IsValidCode(coinCode string) bool {
	ret := _m.Called(coinCode)

	if len(ret) == 0 {
		panic("no return value specified for IsValidCode")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(coinCode)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockCoinValidatorRepository_IsValidCode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsValidCode'
type MockCoinValidatorRepository_IsValidCode_Call struct {
	*mock.Call
}

// IsValidCode is a helper method to define mock.On call
//   - coinCode string
func (_e *MockCoinValidatorRepository_Expecter) IsValidCode(coinCode interface{}) *MockCoinValidatorRepository_IsValidCode_Call {
	return &MockCoinValidatorRepository_IsValidCode_Call{Call: _e.mock.On("IsValidCode", coinCode)}
}

func (_c *MockCoinValidatorRepository_IsValidCode_Call) Run(run func(coinCode string)) *MockCoinValidatorRepository_IsValidCode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCoinValidatorRepository_IsValidCode_Call) Return(_a0 bool) *MockCoinValidatorRepository_IsValidCode_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCoinValidatorRepository_IsValidCode_Call) RunAndReturn(run func(string) bool) *MockCoinValidatorRepository_IsValidCode_Call {
	_c.Call.Return(run)
	return _c
}

// IsValidIdLegacy provides a mock function with given fields: coinIdLegacy
func (_m *MockCoinValidatorRepository) IsValidIdLegacy(coinIdLegacy int32) bool {
	ret := _m.Called(coinIdLegacy)

	if len(ret) == 0 {
		panic("no return value specified for IsValidIdLegacy")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(int32) bool); ok {
		r0 = rf(coinIdLegacy)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockCoinValidatorRepository_IsValidIdLegacy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsValidIdLegacy'
type MockCoinValidatorRepository_IsValidIdLegacy_Call struct {
	*mock.Call
}

// IsValidIdLegacy is a helper method to define mock.On call
//   - coinIdLegacy int32
func (_e *MockCoinValidatorRepository_Expecter) IsValidIdLegacy(coinIdLegacy interface{}) *MockCoinValidatorRepository_IsValidIdLegacy_Call {
	return &MockCoinValidatorRepository_IsValidIdLegacy_Call{Call: _e.mock.On("IsValidIdLegacy", coinIdLegacy)}
}

func (_c *MockCoinValidatorRepository_IsValidIdLegacy_Call) Run(run func(coinIdLegacy int32)) *MockCoinValidatorRepository_IsValidIdLegacy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int32))
	})
	return _c
}

func (_c *MockCoinValidatorRepository_IsValidIdLegacy_Call) Return(_a0 bool) *MockCoinValidatorRepository_IsValidIdLegacy_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCoinValidatorRepository_IsValidIdLegacy_Call) RunAndReturn(run func(int32) bool) *MockCoinValidatorRepository_IsValidIdLegacy_Call {
	_c.Call.Return(run)
	return _c
}

// Serve provides a mock function with given fields: ctx, wg, cacheUpdateInterval
func (_m *MockCoinValidatorRepository) Serve(ctx context.Context, wg *sync.WaitGroup, cacheUpdateInterval time.Duration) {
	_m.Called(ctx, wg, cacheUpdateInterval)
}

// MockCoinValidatorRepository_Serve_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Serve'
type MockCoinValidatorRepository_Serve_Call struct {
	*mock.Call
}

// Serve is a helper method to define mock.On call
//   - ctx context.Context
//   - wg *sync.WaitGroup
//   - cacheUpdateInterval time.Duration
func (_e *MockCoinValidatorRepository_Expecter) Serve(ctx interface{}, wg interface{}, cacheUpdateInterval interface{}) *MockCoinValidatorRepository_Serve_Call {
	return &MockCoinValidatorRepository_Serve_Call{Call: _e.mock.On("Serve", ctx, wg, cacheUpdateInterval)}
}

func (_c *MockCoinValidatorRepository_Serve_Call) Run(run func(ctx context.Context, wg *sync.WaitGroup, cacheUpdateInterval time.Duration)) *MockCoinValidatorRepository_Serve_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*sync.WaitGroup), args[2].(time.Duration))
	})
	return _c
}

func (_c *MockCoinValidatorRepository_Serve_Call) Return() *MockCoinValidatorRepository_Serve_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockCoinValidatorRepository_Serve_Call) RunAndReturn(run func(context.Context, *sync.WaitGroup, time.Duration)) *MockCoinValidatorRepository_Serve_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateMaps provides a mock function with given fields: ctx
func (_m *MockCoinValidatorRepository) UpdateMaps(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMaps")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCoinValidatorRepository_UpdateMaps_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateMaps'
type MockCoinValidatorRepository_UpdateMaps_Call struct {
	*mock.Call
}

// UpdateMaps is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockCoinValidatorRepository_Expecter) UpdateMaps(ctx interface{}) *MockCoinValidatorRepository_UpdateMaps_Call {
	return &MockCoinValidatorRepository_UpdateMaps_Call{Call: _e.mock.On("UpdateMaps", ctx)}
}

func (_c *MockCoinValidatorRepository_UpdateMaps_Call) Run(run func(ctx context.Context)) *MockCoinValidatorRepository_UpdateMaps_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockCoinValidatorRepository_UpdateMaps_Call) Return(_a0 error) *MockCoinValidatorRepository_UpdateMaps_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCoinValidatorRepository_UpdateMaps_Call) RunAndReturn(run func(context.Context) error) *MockCoinValidatorRepository_UpdateMaps_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCoinValidatorRepository creates a new instance of MockCoinValidatorRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCoinValidatorRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCoinValidatorRepository {
	mock := &MockCoinValidatorRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Code generated by mockery v2.46.3. DO NOT EDIT.

package repository

import (
	context "context"

	model "code.emcdtech.com/b2b/swap/model"
	mock "github.com/stretchr/testify/mock"
)

// MockSymbol is an autogenerated mock type for the Symbol type
type MockSymbol struct {
	mock.Mock
}

type MockSymbol_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSymbol) EXPECT() *MockSymbol_Expecter {
	return &MockSymbol_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: ctx, title
func (_m *MockSymbol) Get(ctx context.Context, title string) (*model.Symbol, error) {
	ret := _m.Called(ctx, title)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *model.Symbol
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Symbol, error)); ok {
		return rf(ctx, title)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Symbol); ok {
		r0 = rf(ctx, title)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Symbol)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, title)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSymbol_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockSymbol_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - title string
func (_e *MockSymbol_Expecter) Get(ctx interface{}, title interface{}) *MockSymbol_Get_Call {
	return &MockSymbol_Get_Call{Call: _e.mock.On("Get", ctx, title)}
}

func (_c *MockSymbol_Get_Call) Run(run func(ctx context.Context, title string)) *MockSymbol_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockSymbol_Get_Call) Return(_a0 *model.Symbol, _a1 error) *MockSymbol_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSymbol_Get_Call) RunAndReturn(run func(context.Context, string) (*model.Symbol, error)) *MockSymbol_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetAccuracy provides a mock function with given fields: ctx, symbol
func (_m *MockSymbol) GetAccuracy(ctx context.Context, symbol string) (*model.Accuracy, error) {
	ret := _m.Called(ctx, symbol)

	if len(ret) == 0 {
		panic("no return value specified for GetAccuracy")
	}

	var r0 *model.Accuracy
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Accuracy, error)); ok {
		return rf(ctx, symbol)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Accuracy); ok {
		r0 = rf(ctx, symbol)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Accuracy)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, symbol)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSymbol_GetAccuracy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAccuracy'
type MockSymbol_GetAccuracy_Call struct {
	*mock.Call
}

// GetAccuracy is a helper method to define mock.On call
//   - ctx context.Context
//   - symbol string
func (_e *MockSymbol_Expecter) GetAccuracy(ctx interface{}, symbol interface{}) *MockSymbol_GetAccuracy_Call {
	return &MockSymbol_GetAccuracy_Call{Call: _e.mock.On("GetAccuracy", ctx, symbol)}
}

func (_c *MockSymbol_GetAccuracy_Call) Run(run func(ctx context.Context, symbol string)) *MockSymbol_GetAccuracy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockSymbol_GetAccuracy_Call) Return(_a0 *model.Accuracy, _a1 error) *MockSymbol_GetAccuracy_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSymbol_GetAccuracy_Call) RunAndReturn(run func(context.Context, string) (*model.Accuracy, error)) *MockSymbol_GetAccuracy_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields: ctx
func (_m *MockSymbol) GetAll(ctx context.Context) ([]*model.Symbol, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []*model.Symbol
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*model.Symbol, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*model.Symbol); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Symbol)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSymbol_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type MockSymbol_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockSymbol_Expecter) GetAll(ctx interface{}) *MockSymbol_GetAll_Call {
	return &MockSymbol_GetAll_Call{Call: _e.mock.On("GetAll", ctx)}
}

func (_c *MockSymbol_GetAll_Call) Run(run func(ctx context.Context)) *MockSymbol_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockSymbol_GetAll_Call) Return(_a0 []*model.Symbol, _a1 error) *MockSymbol_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSymbol_GetAll_Call) RunAndReturn(run func(context.Context) ([]*model.Symbol, error)) *MockSymbol_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateAll provides a mock function with given fields: ctx, symbols
func (_m *MockSymbol) UpdateAll(ctx context.Context, symbols map[string]*model.Symbol) error {
	ret := _m.Called(ctx, symbols)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAll")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]*model.Symbol) error); ok {
		r0 = rf(ctx, symbols)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSymbol_UpdateAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateAll'
type MockSymbol_UpdateAll_Call struct {
	*mock.Call
}

// UpdateAll is a helper method to define mock.On call
//   - ctx context.Context
//   - symbols map[string]*model.Symbol
func (_e *MockSymbol_Expecter) UpdateAll(ctx interface{}, symbols interface{}) *MockSymbol_UpdateAll_Call {
	return &MockSymbol_UpdateAll_Call{Call: _e.mock.On("UpdateAll", ctx, symbols)}
}

func (_c *MockSymbol_UpdateAll_Call) Run(run func(ctx context.Context, symbols map[string]*model.Symbol)) *MockSymbol_UpdateAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(map[string]*model.Symbol))
	})
	return _c
}

func (_c *MockSymbol_UpdateAll_Call) Return(_a0 error) *MockSymbol_UpdateAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSymbol_UpdateAll_Call) RunAndReturn(run func(context.Context, map[string]*model.Symbol) error) *MockSymbol_UpdateAll_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSymbol creates a new instance of MockSymbol. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSymbol(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSymbol {
	mock := &MockSymbol{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

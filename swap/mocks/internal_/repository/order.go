// Code generated by mockery v2.46.3. DO NOT EDIT.

package repository

import (
	context "context"

	model "code.emcdtech.com/b2b/swap/model"
	mock "github.com/stretchr/testify/mock"
)

// MockOrder is an autogenerated mock type for the Order type
type MockOrder struct {
	mock.Mock
}

type MockOrder_Expecter struct {
	mock *mock.Mock
}

func (_m *MockOrder) EXPECT() *MockOrder_Expecter {
	return &MockOrder_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: ctx, order
func (_m *MockOrder) Add(ctx context.Context, order *model.Order) error {
	ret := _m.Called(ctx, order)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Order) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockOrder_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type MockOrder_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - ctx context.Context
//   - order *model.Order
func (_e *MockOrder_Expecter) Add(ctx interface{}, order interface{}) *MockOrder_Add_Call {
	return &MockOrder_Add_Call{Call: _e.mock.On("Add", ctx, order)}
}

func (_c *MockOrder_Add_Call) Run(run func(ctx context.Context, order *model.Order)) *MockOrder_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Order))
	})
	return _c
}

func (_c *MockOrder_Add_Call) Return(_a0 error) *MockOrder_Add_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockOrder_Add_Call) RunAndReturn(run func(context.Context, *model.Order) error) *MockOrder_Add_Call {
	_c.Call.Return(run)
	return _c
}

// Find provides a mock function with given fields: ctx, filter
func (_m *MockOrder) Find(ctx context.Context, filter *model.OrderFilter) (model.Orders, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 model.Orders
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.OrderFilter) (model.Orders, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.OrderFilter) model.Orders); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.Orders)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.OrderFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrder_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type MockOrder_Find_Call struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - ctx context.Context
//   - filter *model.OrderFilter
func (_e *MockOrder_Expecter) Find(ctx interface{}, filter interface{}) *MockOrder_Find_Call {
	return &MockOrder_Find_Call{Call: _e.mock.On("Find", ctx, filter)}
}

func (_c *MockOrder_Find_Call) Run(run func(ctx context.Context, filter *model.OrderFilter)) *MockOrder_Find_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.OrderFilter))
	})
	return _c
}

func (_c *MockOrder_Find_Call) Return(_a0 model.Orders, _a1 error) *MockOrder_Find_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrder_Find_Call) RunAndReturn(run func(context.Context, *model.OrderFilter) (model.Orders, error)) *MockOrder_Find_Call {
	_c.Call.Return(run)
	return _c
}

// FindOne provides a mock function with given fields: ctx, filter
func (_m *MockOrder) FindOne(ctx context.Context, filter *model.OrderFilter) (*model.Order, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for FindOne")
	}

	var r0 *model.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.OrderFilter) (*model.Order, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.OrderFilter) *model.Order); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.OrderFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockOrder_FindOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindOne'
type MockOrder_FindOne_Call struct {
	*mock.Call
}

// FindOne is a helper method to define mock.On call
//   - ctx context.Context
//   - filter *model.OrderFilter
func (_e *MockOrder_Expecter) FindOne(ctx interface{}, filter interface{}) *MockOrder_FindOne_Call {
	return &MockOrder_FindOne_Call{Call: _e.mock.On("FindOne", ctx, filter)}
}

func (_c *MockOrder_FindOne_Call) Run(run func(ctx context.Context, filter *model.OrderFilter)) *MockOrder_FindOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.OrderFilter))
	})
	return _c
}

func (_c *MockOrder_FindOne_Call) Return(_a0 *model.Order, _a1 error) *MockOrder_FindOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockOrder_FindOne_Call) RunAndReturn(run func(context.Context, *model.OrderFilter) (*model.Order, error)) *MockOrder_FindOne_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, order, filter, partial
func (_m *MockOrder) Update(ctx context.Context, order *model.Order, filter *model.OrderFilter, partial *model.OrderPartial) error {
	ret := _m.Called(ctx, order, filter, partial)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Order, *model.OrderFilter, *model.OrderPartial) error); ok {
		r0 = rf(ctx, order, filter, partial)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockOrder_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockOrder_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - order *model.Order
//   - filter *model.OrderFilter
//   - partial *model.OrderPartial
func (_e *MockOrder_Expecter) Update(ctx interface{}, order interface{}, filter interface{}, partial interface{}) *MockOrder_Update_Call {
	return &MockOrder_Update_Call{Call: _e.mock.On("Update", ctx, order, filter, partial)}
}

func (_c *MockOrder_Update_Call) Run(run func(ctx context.Context, order *model.Order, filter *model.OrderFilter, partial *model.OrderPartial)) *MockOrder_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Order), args[2].(*model.OrderFilter), args[3].(*model.OrderPartial))
	})
	return _c
}

func (_c *MockOrder_Update_Call) Return(_a0 error) *MockOrder_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockOrder_Update_Call) RunAndReturn(run func(context.Context, *model.Order, *model.OrderFilter, *model.OrderPartial) error) *MockOrder_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockOrder creates a new instance of MockOrder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockOrder(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockOrder {
	mock := &MockOrder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

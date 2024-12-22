// Code generated by mockery v2.46.3. DO NOT EDIT.

package repository

import (
	model "code.emcdtech.com/b2b/swap/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockStatusCache is an autogenerated mock type for the StatusCache type
type MockStatusCache struct {
	mock.Mock
}

type MockStatusCache_Expecter struct {
	mock *mock.Mock
}

func (_m *MockStatusCache) EXPECT() *MockStatusCache_Expecter {
	return &MockStatusCache_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: swapID, status
func (_m *MockStatusCache) Add(swapID uuid.UUID, status model.PublicStatus) {
	_m.Called(swapID, status)
}

// MockStatusCache_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type MockStatusCache_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - swapID uuid.UUID
//   - status model.PublicStatus
func (_e *MockStatusCache_Expecter) Add(swapID interface{}, status interface{}) *MockStatusCache_Add_Call {
	return &MockStatusCache_Add_Call{Call: _e.mock.On("Add", swapID, status)}
}

func (_c *MockStatusCache_Add_Call) Run(run func(swapID uuid.UUID, status model.PublicStatus)) *MockStatusCache_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID), args[1].(model.PublicStatus))
	})
	return _c
}

func (_c *MockStatusCache_Add_Call) Return() *MockStatusCache_Add_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockStatusCache_Add_Call) RunAndReturn(run func(uuid.UUID, model.PublicStatus)) *MockStatusCache_Add_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: swapID
func (_m *MockStatusCache) Delete(swapID uuid.UUID) {
	_m.Called(swapID)
}

// MockStatusCache_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockStatusCache_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - swapID uuid.UUID
func (_e *MockStatusCache_Expecter) Delete(swapID interface{}) *MockStatusCache_Delete_Call {
	return &MockStatusCache_Delete_Call{Call: _e.mock.On("Delete", swapID)}
}

func (_c *MockStatusCache_Delete_Call) Run(run func(swapID uuid.UUID)) *MockStatusCache_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID))
	})
	return _c
}

func (_c *MockStatusCache_Delete_Call) Return() *MockStatusCache_Delete_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockStatusCache_Delete_Call) RunAndReturn(run func(uuid.UUID)) *MockStatusCache_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: swapID
func (_m *MockStatusCache) Get(swapID uuid.UUID) model.PublicStatus {
	ret := _m.Called(swapID)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 model.PublicStatus
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.PublicStatus); ok {
		r0 = rf(swapID)
	} else {
		r0 = ret.Get(0).(model.PublicStatus)
	}

	return r0
}

// MockStatusCache_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockStatusCache_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - swapID uuid.UUID
func (_e *MockStatusCache_Expecter) Get(swapID interface{}) *MockStatusCache_Get_Call {
	return &MockStatusCache_Get_Call{Call: _e.mock.On("Get", swapID)}
}

func (_c *MockStatusCache_Get_Call) Run(run func(swapID uuid.UUID)) *MockStatusCache_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID))
	})
	return _c
}

func (_c *MockStatusCache_Get_Call) Return(_a0 model.PublicStatus) *MockStatusCache_Get_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockStatusCache_Get_Call) RunAndReturn(run func(uuid.UUID) model.PublicStatus) *MockStatusCache_Get_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockStatusCache creates a new instance of MockStatusCache. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockStatusCache(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockStatusCache {
	mock := &MockStatusCache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

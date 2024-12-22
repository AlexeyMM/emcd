// Code generated by mockery v2.46.3. DO NOT EDIT.

package swapCoin

import mock "github.com/stretchr/testify/mock"

// MockUnsafeSwapCoinServiceServer is an autogenerated mock type for the UnsafeSwapCoinServiceServer type
type MockUnsafeSwapCoinServiceServer struct {
	mock.Mock
}

type MockUnsafeSwapCoinServiceServer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUnsafeSwapCoinServiceServer) EXPECT() *MockUnsafeSwapCoinServiceServer_Expecter {
	return &MockUnsafeSwapCoinServiceServer_Expecter{mock: &_m.Mock}
}

// mustEmbedUnimplementedSwapCoinServiceServer provides a mock function with given fields:
func (_m *MockUnsafeSwapCoinServiceServer) mustEmbedUnimplementedSwapCoinServiceServer() {
	_m.Called()
}

// MockUnsafeSwapCoinServiceServer_mustEmbedUnimplementedSwapCoinServiceServer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'mustEmbedUnimplementedSwapCoinServiceServer'
type MockUnsafeSwapCoinServiceServer_mustEmbedUnimplementedSwapCoinServiceServer_Call struct {
	*mock.Call
}

// mustEmbedUnimplementedSwapCoinServiceServer is a helper method to define mock.On call
func (_e *MockUnsafeSwapCoinServiceServer_Expecter) mustEmbedUnimplementedSwapCoinServiceServer() *MockUnsafeSwapCoinServiceServer_mustEmbedUnimplementedSwapCoinServiceServer_Call {
	return &MockUnsafeSwapCoinServiceServer_mustEmbedUnimplementedSwapCoinServiceServer_Call{Call: _e.mock.On("mustEmbedUnimplementedSwapCoinServiceServer")}
}

func (_c *MockUnsafeSwapCoinServiceServer_mustEmbedUnimplementedSwapCoinServiceServer_Call) Run(run func()) *MockUnsafeSwapCoinServiceServer_mustEmbedUnimplementedSwapCoinServiceServer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockUnsafeSwapCoinServiceServer_mustEmbedUnimplementedSwapCoinServiceServer_Call) Return() *MockUnsafeSwapCoinServiceServer_mustEmbedUnimplementedSwapCoinServiceServer_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockUnsafeSwapCoinServiceServer_mustEmbedUnimplementedSwapCoinServiceServer_Call) RunAndReturn(run func()) *MockUnsafeSwapCoinServiceServer_mustEmbedUnimplementedSwapCoinServiceServer_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUnsafeSwapCoinServiceServer creates a new instance of MockUnsafeSwapCoinServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUnsafeSwapCoinServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUnsafeSwapCoinServiceServer {
	mock := &MockUnsafeSwapCoinServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
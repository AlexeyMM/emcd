// Code generated by mockery v2.43.2. DO NOT EDIT.

package referral

import mock "github.com/stretchr/testify/mock"

// MockUnsafeAccountingReferralServiceServer is an autogenerated mock type for the UnsafeAccountingReferralServiceServer type
type MockUnsafeAccountingReferralServiceServer struct {
	mock.Mock
}

type MockUnsafeAccountingReferralServiceServer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUnsafeAccountingReferralServiceServer) EXPECT() *MockUnsafeAccountingReferralServiceServer_Expecter {
	return &MockUnsafeAccountingReferralServiceServer_Expecter{mock: &_m.Mock}
}

// mustEmbedUnimplementedAccountingReferralServiceServer provides a mock function with given fields:
func (_m *MockUnsafeAccountingReferralServiceServer) mustEmbedUnimplementedAccountingReferralServiceServer() {
	_m.Called()
}

// MockUnsafeAccountingReferralServiceServer_mustEmbedUnimplementedAccountingReferralServiceServer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'mustEmbedUnimplementedAccountingReferralServiceServer'
type MockUnsafeAccountingReferralServiceServer_mustEmbedUnimplementedAccountingReferralServiceServer_Call struct {
	*mock.Call
}

// mustEmbedUnimplementedAccountingReferralServiceServer is a helper method to define mock.On call
func (_e *MockUnsafeAccountingReferralServiceServer_Expecter) mustEmbedUnimplementedAccountingReferralServiceServer() *MockUnsafeAccountingReferralServiceServer_mustEmbedUnimplementedAccountingReferralServiceServer_Call {
	return &MockUnsafeAccountingReferralServiceServer_mustEmbedUnimplementedAccountingReferralServiceServer_Call{Call: _e.mock.On("mustEmbedUnimplementedAccountingReferralServiceServer")}
}

func (_c *MockUnsafeAccountingReferralServiceServer_mustEmbedUnimplementedAccountingReferralServiceServer_Call) Run(run func()) *MockUnsafeAccountingReferralServiceServer_mustEmbedUnimplementedAccountingReferralServiceServer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockUnsafeAccountingReferralServiceServer_mustEmbedUnimplementedAccountingReferralServiceServer_Call) Return() *MockUnsafeAccountingReferralServiceServer_mustEmbedUnimplementedAccountingReferralServiceServer_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockUnsafeAccountingReferralServiceServer_mustEmbedUnimplementedAccountingReferralServiceServer_Call) RunAndReturn(run func()) *MockUnsafeAccountingReferralServiceServer_mustEmbedUnimplementedAccountingReferralServiceServer_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUnsafeAccountingReferralServiceServer creates a new instance of MockUnsafeAccountingReferralServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUnsafeAccountingReferralServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUnsafeAccountingReferralServiceServer {
	mock := &MockUnsafeAccountingReferralServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

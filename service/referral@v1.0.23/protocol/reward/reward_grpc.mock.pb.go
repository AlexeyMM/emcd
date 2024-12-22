// Code generated by MockGen. DO NOT EDIT.
// Source: protocol/reward/reward_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -source protocol/reward/reward_grpc.pb.go -destination protocol/reward/reward_grpc.mock.pb.go -package reward
//

// Package reward is a generated GoMock package.
package reward

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockRewardServiceClient is a mock of RewardServiceClient interface.
type MockRewardServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockRewardServiceClientMockRecorder
}

// MockRewardServiceClientMockRecorder is the mock recorder for MockRewardServiceClient.
type MockRewardServiceClientMockRecorder struct {
	mock *MockRewardServiceClient
}

// NewMockRewardServiceClient creates a new mock instance.
func NewMockRewardServiceClient(ctrl *gomock.Controller) *MockRewardServiceClient {
	mock := &MockRewardServiceClient{ctrl: ctrl}
	mock.recorder = &MockRewardServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRewardServiceClient) EXPECT() *MockRewardServiceClientMockRecorder {
	return m.recorder
}

// Calculate mocks base method.
func (m *MockRewardServiceClient) Calculate(ctx context.Context, in *CalculateRequest, opts ...grpc.CallOption) (*CalculateResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Calculate", varargs...)
	ret0, _ := ret[0].(*CalculateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Calculate indicates an expected call of Calculate.
func (mr *MockRewardServiceClientMockRecorder) Calculate(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Calculate", reflect.TypeOf((*MockRewardServiceClient)(nil).Calculate), varargs...)
}

// UpdateWithMultiplier mocks base method.
func (m *MockRewardServiceClient) UpdateWithMultiplier(ctx context.Context, in *UpdateWithMultiplierRequest, opts ...grpc.CallOption) (*UpdateWithMultiplierResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateWithMultiplier", varargs...)
	ret0, _ := ret[0].(*UpdateWithMultiplierResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateWithMultiplier indicates an expected call of UpdateWithMultiplier.
func (mr *MockRewardServiceClientMockRecorder) UpdateWithMultiplier(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWithMultiplier", reflect.TypeOf((*MockRewardServiceClient)(nil).UpdateWithMultiplier), varargs...)
}

// MockRewardServiceServer is a mock of RewardServiceServer interface.
type MockRewardServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockRewardServiceServerMockRecorder
}

// MockRewardServiceServerMockRecorder is the mock recorder for MockRewardServiceServer.
type MockRewardServiceServerMockRecorder struct {
	mock *MockRewardServiceServer
}

// NewMockRewardServiceServer creates a new mock instance.
func NewMockRewardServiceServer(ctrl *gomock.Controller) *MockRewardServiceServer {
	mock := &MockRewardServiceServer{ctrl: ctrl}
	mock.recorder = &MockRewardServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRewardServiceServer) EXPECT() *MockRewardServiceServerMockRecorder {
	return m.recorder
}

// Calculate mocks base method.
func (m *MockRewardServiceServer) Calculate(arg0 context.Context, arg1 *CalculateRequest) (*CalculateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Calculate", arg0, arg1)
	ret0, _ := ret[0].(*CalculateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Calculate indicates an expected call of Calculate.
func (mr *MockRewardServiceServerMockRecorder) Calculate(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Calculate", reflect.TypeOf((*MockRewardServiceServer)(nil).Calculate), arg0, arg1)
}

// UpdateWithMultiplier mocks base method.
func (m *MockRewardServiceServer) UpdateWithMultiplier(arg0 context.Context, arg1 *UpdateWithMultiplierRequest) (*UpdateWithMultiplierResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWithMultiplier", arg0, arg1)
	ret0, _ := ret[0].(*UpdateWithMultiplierResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateWithMultiplier indicates an expected call of UpdateWithMultiplier.
func (mr *MockRewardServiceServerMockRecorder) UpdateWithMultiplier(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWithMultiplier", reflect.TypeOf((*MockRewardServiceServer)(nil).UpdateWithMultiplier), arg0, arg1)
}

// mustEmbedUnimplementedRewardServiceServer mocks base method.
func (m *MockRewardServiceServer) mustEmbedUnimplementedRewardServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedRewardServiceServer")
}

// mustEmbedUnimplementedRewardServiceServer indicates an expected call of mustEmbedUnimplementedRewardServiceServer.
func (mr *MockRewardServiceServerMockRecorder) mustEmbedUnimplementedRewardServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedRewardServiceServer", reflect.TypeOf((*MockRewardServiceServer)(nil).mustEmbedUnimplementedRewardServiceServer))
}

// MockUnsafeRewardServiceServer is a mock of UnsafeRewardServiceServer interface.
type MockUnsafeRewardServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeRewardServiceServerMockRecorder
}

// MockUnsafeRewardServiceServerMockRecorder is the mock recorder for MockUnsafeRewardServiceServer.
type MockUnsafeRewardServiceServerMockRecorder struct {
	mock *MockUnsafeRewardServiceServer
}

// NewMockUnsafeRewardServiceServer creates a new mock instance.
func NewMockUnsafeRewardServiceServer(ctrl *gomock.Controller) *MockUnsafeRewardServiceServer {
	mock := &MockUnsafeRewardServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeRewardServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeRewardServiceServer) EXPECT() *MockUnsafeRewardServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedRewardServiceServer mocks base method.
func (m *MockUnsafeRewardServiceServer) mustEmbedUnimplementedRewardServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedRewardServiceServer")
}

// mustEmbedUnimplementedRewardServiceServer indicates an expected call of mustEmbedUnimplementedRewardServiceServer.
func (mr *MockUnsafeRewardServiceServerMockRecorder) mustEmbedUnimplementedRewardServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedRewardServiceServer", reflect.TypeOf((*MockUnsafeRewardServiceServer)(nil).mustEmbedUnimplementedRewardServiceServer))
}

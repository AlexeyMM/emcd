// Code generated by MockGen. DO NOT EDIT.
// Source: email_messages.go
//
// Generated by this command:
//
//	mockgen -source=email_messages.go -destination=./mocks/email_messages.go -package=mockstore
//

// Package mockstore is a generated GoMock package.
package mockstore

import (
	context "context"
	reflect "reflect"

	model "code.emcdtech.com/emcd/service/email/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockEmailMessages is a mock of EmailMessages interface.
type MockEmailMessages struct {
	ctrl     *gomock.Controller
	recorder *MockEmailMessagesMockRecorder
	isgomock struct{}
}

// MockEmailMessagesMockRecorder is the mock recorder for MockEmailMessages.
type MockEmailMessagesMockRecorder struct {
	mock *MockEmailMessages
}

// NewMockEmailMessages creates a new mock instance.
func NewMockEmailMessages(ctrl *gomock.Controller) *MockEmailMessages {
	mock := &MockEmailMessages{ctrl: ctrl}
	mock.recorder = &MockEmailMessagesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailMessages) EXPECT() *MockEmailMessagesMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockEmailMessages) Create(ctx context.Context, em *model.EmailMessageEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, em)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockEmailMessagesMockRecorder) Create(ctx, em any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockEmailMessages)(nil).Create), ctx, em)
}

// ListMessages mocks base method.
func (m *MockEmailMessages) ListMessages(ctx context.Context, email, eventType *string, skip, take int32) ([]*model.EmailMessageEvent, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListMessages", ctx, email, eventType, skip, take)
	ret0, _ := ret[0].([]*model.EmailMessageEvent)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListMessages indicates an expected call of ListMessages.
func (mr *MockEmailMessagesMockRecorder) ListMessages(ctx, email, eventType, skip, take any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListMessages", reflect.TypeOf((*MockEmailMessages)(nil).ListMessages), ctx, email, eventType, skip, take)
}
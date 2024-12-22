// Code generated by MockGen. DO NOT EDIT.
// Source: template.go
//
// Generated by this command:
//
//	mockgen -source=template.go -destination=./mocks/templates.go -package=mockstore
//

// Package mockstore is a generated GoMock package.
package mockstore

import (
	context "context"
	reflect "reflect"

	model "code.emcdtech.com/emcd/service/email/internal/model"
	repository "code.emcdtech.com/emcd/service/email/internal/repository"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockTemplate is a mock of Template interface.
type MockTemplate struct {
	ctrl     *gomock.Controller
	recorder *MockTemplateMockRecorder
	isgomock struct{}
}

// MockTemplateMockRecorder is the mock recorder for MockTemplate.
type MockTemplateMockRecorder struct {
	mock *MockTemplate
}

// NewMockTemplate creates a new mock instance.
func NewMockTemplate(ctrl *gomock.Controller) *MockTemplate {
	mock := &MockTemplate{ctrl: ctrl}
	mock.recorder = &MockTemplateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTemplate) EXPECT() *MockTemplateMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTemplate) Create(ctx context.Context, template model.Template) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, template)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockTemplateMockRecorder) Create(ctx, template any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTemplate)(nil).Create), ctx, template)
}

// Get mocks base method.
func (m *MockTemplate) Get(ctx context.Context, whitelabelID uuid.UUID, language string, _type model.CodeTemplate) (model.Template, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, whitelabelID, language, _type)
	ret0, _ := ret[0].(model.Template)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockTemplateMockRecorder) Get(ctx, whitelabelID, language, _type any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTemplate)(nil).Get), ctx, whitelabelID, language, _type)
}

// List mocks base method.
func (m *MockTemplate) List(ctx context.Context, pagination repository.Pagination) ([]model.Template, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, pagination)
	ret0, _ := ret[0].([]model.Template)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockTemplateMockRecorder) List(ctx, pagination any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockTemplate)(nil).List), ctx, pagination)
}

// Update mocks base method.
func (m *MockTemplate) Update(ctx context.Context, template model.Template) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, template)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTemplateMockRecorder) Update(ctx, template any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTemplate)(nil).Update), ctx, template)
}
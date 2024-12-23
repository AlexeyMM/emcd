// Code generated by mockery v2.43.2. DO NOT EDIT.

package service

import (
	context "context"

	enum "code.emcdtech.com/emcd/service/accounting/model/enum"
	mock "github.com/stretchr/testify/mock"

	model "code.emcdtech.com/emcd/service/accounting/model"

	uuid "github.com/google/uuid"
)

// MockUserAccountService is an autogenerated mock type for the UserAccountService type
type MockUserAccountService struct {
	mock.Mock
}

type MockUserAccountService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUserAccountService) EXPECT() *MockUserAccountService_Expecter {
	return &MockUserAccountService_Expecter{mock: &_m.Mock}
}

// CreateUserAccounts provides a mock function with given fields: ctx, userId, userIdNew, userAccounts
func (_m *MockUserAccountService) CreateUserAccounts(ctx context.Context, userId int32, userIdNew uuid.UUID, userAccounts model.UserAccounts) (model.UserAccounts, error) {
	ret := _m.Called(ctx, userId, userIdNew, userAccounts)

	if len(ret) == 0 {
		panic("no return value specified for CreateUserAccounts")
	}

	var r0 model.UserAccounts
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int32, uuid.UUID, model.UserAccounts) (model.UserAccounts, error)); ok {
		return rf(ctx, userId, userIdNew, userAccounts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int32, uuid.UUID, model.UserAccounts) model.UserAccounts); ok {
		r0 = rf(ctx, userId, userIdNew, userAccounts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.UserAccounts)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int32, uuid.UUID, model.UserAccounts) error); ok {
		r1 = rf(ctx, userId, userIdNew, userAccounts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserAccountService_CreateUserAccounts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUserAccounts'
type MockUserAccountService_CreateUserAccounts_Call struct {
	*mock.Call
}

// CreateUserAccounts is a helper method to define mock.On call
//   - ctx context.Context
//   - userId int32
//   - userIdNew uuid.UUID
//   - userAccounts model.UserAccounts
func (_e *MockUserAccountService_Expecter) CreateUserAccounts(ctx interface{}, userId interface{}, userIdNew interface{}, userAccounts interface{}) *MockUserAccountService_CreateUserAccounts_Call {
	return &MockUserAccountService_CreateUserAccounts_Call{Call: _e.mock.On("CreateUserAccounts", ctx, userId, userIdNew, userAccounts)}
}

func (_c *MockUserAccountService_CreateUserAccounts_Call) Run(run func(ctx context.Context, userId int32, userIdNew uuid.UUID, userAccounts model.UserAccounts)) *MockUserAccountService_CreateUserAccounts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int32), args[2].(uuid.UUID), args[3].(model.UserAccounts))
	})
	return _c
}

func (_c *MockUserAccountService_CreateUserAccounts_Call) Return(_a0 model.UserAccounts, _a1 error) *MockUserAccountService_CreateUserAccounts_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserAccountService_CreateUserAccounts_Call) RunAndReturn(run func(context.Context, int32, uuid.UUID, model.UserAccounts) (model.UserAccounts, error)) *MockUserAccountService_CreateUserAccounts_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrCreateUserAccount provides a mock function with given fields: ctx, userAccount
func (_m *MockUserAccountService) GetOrCreateUserAccount(ctx context.Context, userAccount *model.UserAccount) (*model.UserAccount, error) {
	ret := _m.Called(ctx, userAccount)

	if len(ret) == 0 {
		panic("no return value specified for GetOrCreateUserAccount")
	}

	var r0 *model.UserAccount
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserAccount) (*model.UserAccount, error)); ok {
		return rf(ctx, userAccount)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserAccount) *model.UserAccount); ok {
		r0 = rf(ctx, userAccount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserAccount)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.UserAccount) error); ok {
		r1 = rf(ctx, userAccount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserAccountService_GetOrCreateUserAccount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrCreateUserAccount'
type MockUserAccountService_GetOrCreateUserAccount_Call struct {
	*mock.Call
}

// GetOrCreateUserAccount is a helper method to define mock.On call
//   - ctx context.Context
//   - userAccount *model.UserAccount
func (_e *MockUserAccountService_Expecter) GetOrCreateUserAccount(ctx interface{}, userAccount interface{}) *MockUserAccountService_GetOrCreateUserAccount_Call {
	return &MockUserAccountService_GetOrCreateUserAccount_Call{Call: _e.mock.On("GetOrCreateUserAccount", ctx, userAccount)}
}

func (_c *MockUserAccountService_GetOrCreateUserAccount_Call) Run(run func(ctx context.Context, userAccount *model.UserAccount)) *MockUserAccountService_GetOrCreateUserAccount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.UserAccount))
	})
	return _c
}

func (_c *MockUserAccountService_GetOrCreateUserAccount_Call) Return(_a0 *model.UserAccount, _a1 error) *MockUserAccountService_GetOrCreateUserAccount_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserAccountService_GetOrCreateUserAccount_Call) RunAndReturn(run func(context.Context, *model.UserAccount) (*model.UserAccount, error)) *MockUserAccountService_GetOrCreateUserAccount_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserAccountByConstraint provides a mock function with given fields: ctx, userIdNew, coinIdNew, userAccountId
func (_m *MockUserAccountService) GetUserAccountByConstraint(ctx context.Context, userIdNew uuid.UUID, coinIdNew string, userAccountId enum.AccountTypeId) (*model.UserAccount, error) {
	ret := _m.Called(ctx, userIdNew, coinIdNew, userAccountId)

	if len(ret) == 0 {
		panic("no return value specified for GetUserAccountByConstraint")
	}

	var r0 *model.UserAccount
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string, enum.AccountTypeId) (*model.UserAccount, error)); ok {
		return rf(ctx, userIdNew, coinIdNew, userAccountId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string, enum.AccountTypeId) *model.UserAccount); ok {
		r0 = rf(ctx, userIdNew, coinIdNew, userAccountId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserAccount)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, string, enum.AccountTypeId) error); ok {
		r1 = rf(ctx, userIdNew, coinIdNew, userAccountId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserAccountService_GetUserAccountByConstraint_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserAccountByConstraint'
type MockUserAccountService_GetUserAccountByConstraint_Call struct {
	*mock.Call
}

// GetUserAccountByConstraint is a helper method to define mock.On call
//   - ctx context.Context
//   - userIdNew uuid.UUID
//   - coinIdNew string
//   - userAccountId enum.AccountTypeId
func (_e *MockUserAccountService_Expecter) GetUserAccountByConstraint(ctx interface{}, userIdNew interface{}, coinIdNew interface{}, userAccountId interface{}) *MockUserAccountService_GetUserAccountByConstraint_Call {
	return &MockUserAccountService_GetUserAccountByConstraint_Call{Call: _e.mock.On("GetUserAccountByConstraint", ctx, userIdNew, coinIdNew, userAccountId)}
}

func (_c *MockUserAccountService_GetUserAccountByConstraint_Call) Run(run func(ctx context.Context, userIdNew uuid.UUID, coinIdNew string, userAccountId enum.AccountTypeId)) *MockUserAccountService_GetUserAccountByConstraint_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(string), args[3].(enum.AccountTypeId))
	})
	return _c
}

func (_c *MockUserAccountService_GetUserAccountByConstraint_Call) Return(_a0 *model.UserAccount, _a1 error) *MockUserAccountService_GetUserAccountByConstraint_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserAccountService_GetUserAccountByConstraint_Call) RunAndReturn(run func(context.Context, uuid.UUID, string, enum.AccountTypeId) (*model.UserAccount, error)) *MockUserAccountService_GetUserAccountByConstraint_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserAccountById provides a mock function with given fields: ctx, id
func (_m *MockUserAccountService) GetUserAccountById(ctx context.Context, id int32) (*model.UserAccount, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetUserAccountById")
	}

	var r0 *model.UserAccount
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int32) (*model.UserAccount, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int32) *model.UserAccount); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserAccount)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int32) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserAccountService_GetUserAccountById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserAccountById'
type MockUserAccountService_GetUserAccountById_Call struct {
	*mock.Call
}

// GetUserAccountById is a helper method to define mock.On call
//   - ctx context.Context
//   - id int32
func (_e *MockUserAccountService_Expecter) GetUserAccountById(ctx interface{}, id interface{}) *MockUserAccountService_GetUserAccountById_Call {
	return &MockUserAccountService_GetUserAccountById_Call{Call: _e.mock.On("GetUserAccountById", ctx, id)}
}

func (_c *MockUserAccountService_GetUserAccountById_Call) Run(run func(ctx context.Context, id int32)) *MockUserAccountService_GetUserAccountById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int32))
	})
	return _c
}

func (_c *MockUserAccountService_GetUserAccountById_Call) Return(_a0 *model.UserAccount, _a1 error) *MockUserAccountService_GetUserAccountById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserAccountService_GetUserAccountById_Call) RunAndReturn(run func(context.Context, int32) (*model.UserAccount, error)) *MockUserAccountService_GetUserAccountById_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserAccountsByFilter provides a mock function with given fields: ctx, filter
func (_m *MockUserAccountService) GetUserAccountsByFilter(ctx context.Context, filter *model.UserAccountFilter) (*uint64, model.UserAccounts, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for GetUserAccountsByFilter")
	}

	var r0 *uint64
	var r1 model.UserAccounts
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserAccountFilter) (*uint64, model.UserAccounts, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserAccountFilter) *uint64); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uint64)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.UserAccountFilter) model.UserAccounts); ok {
		r1 = rf(ctx, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(model.UserAccounts)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, *model.UserAccountFilter) error); ok {
		r2 = rf(ctx, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockUserAccountService_GetUserAccountsByFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserAccountsByFilter'
type MockUserAccountService_GetUserAccountsByFilter_Call struct {
	*mock.Call
}

// GetUserAccountsByFilter is a helper method to define mock.On call
//   - ctx context.Context
//   - filter *model.UserAccountFilter
func (_e *MockUserAccountService_Expecter) GetUserAccountsByFilter(ctx interface{}, filter interface{}) *MockUserAccountService_GetUserAccountsByFilter_Call {
	return &MockUserAccountService_GetUserAccountsByFilter_Call{Call: _e.mock.On("GetUserAccountsByFilter", ctx, filter)}
}

func (_c *MockUserAccountService_GetUserAccountsByFilter_Call) Run(run func(ctx context.Context, filter *model.UserAccountFilter)) *MockUserAccountService_GetUserAccountsByFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.UserAccountFilter))
	})
	return _c
}

func (_c *MockUserAccountService_GetUserAccountsByFilter_Call) Return(_a0 *uint64, _a1 model.UserAccounts, _a2 error) *MockUserAccountService_GetUserAccountsByFilter_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockUserAccountService_GetUserAccountsByFilter_Call) RunAndReturn(run func(context.Context, *model.UserAccountFilter) (*uint64, model.UserAccounts, error)) *MockUserAccountService_GetUserAccountsByFilter_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserAccountsByUuid provides a mock function with given fields: ctx, userUuid
func (_m *MockUserAccountService) GetUserAccountsByUuid(ctx context.Context, userUuid uuid.UUID) (model.UserAccounts, error) {
	ret := _m.Called(ctx, userUuid)

	if len(ret) == 0 {
		panic("no return value specified for GetUserAccountsByUuid")
	}

	var r0 model.UserAccounts
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (model.UserAccounts, error)); ok {
		return rf(ctx, userUuid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) model.UserAccounts); ok {
		r0 = rf(ctx, userUuid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.UserAccounts)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, userUuid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserAccountService_GetUserAccountsByUuid_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserAccountsByUuid'
type MockUserAccountService_GetUserAccountsByUuid_Call struct {
	*mock.Call
}

// GetUserAccountsByUuid is a helper method to define mock.On call
//   - ctx context.Context
//   - userUuid uuid.UUID
func (_e *MockUserAccountService_Expecter) GetUserAccountsByUuid(ctx interface{}, userUuid interface{}) *MockUserAccountService_GetUserAccountsByUuid_Call {
	return &MockUserAccountService_GetUserAccountsByUuid_Call{Call: _e.mock.On("GetUserAccountsByUuid", ctx, userUuid)}
}

func (_c *MockUserAccountService_GetUserAccountsByUuid_Call) Run(run func(ctx context.Context, userUuid uuid.UUID)) *MockUserAccountService_GetUserAccountsByUuid_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockUserAccountService_GetUserAccountsByUuid_Call) Return(_a0 model.UserAccounts, _a1 error) *MockUserAccountService_GetUserAccountsByUuid_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserAccountService_GetUserAccountsByUuid_Call) RunAndReturn(run func(context.Context, uuid.UUID) (model.UserAccounts, error)) *MockUserAccountService_GetUserAccountsByUuid_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUserAccountService creates a new instance of MockUserAccountService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserAccountService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserAccountService {
	mock := &MockUserAccountService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

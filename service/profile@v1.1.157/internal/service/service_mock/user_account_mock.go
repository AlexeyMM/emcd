// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	context "context"

	model "code.emcdtech.com/emcd/service/accounting/model"
	mock "github.com/stretchr/testify/mock"

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

// CreateUserAccount provides a mock function with given fields: ctx, userId, userIdNew, coinId, userAccountTypeId, minpay
func (_m *MockUserAccountService) CreateUserAccount(ctx context.Context, userId int32, userIdNew uuid.UUID, coinId int32, userAccountTypeId int32, minpay float64) (int, error) {
	ret := _m.Called(ctx, userId, userIdNew, coinId, userAccountTypeId, minpay)

	if len(ret) == 0 {
		panic("no return value specified for CreateUserAccount")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int32, uuid.UUID, int32, int32, float64) (int, error)); ok {
		return rf(ctx, userId, userIdNew, coinId, userAccountTypeId, minpay)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int32, uuid.UUID, int32, int32, float64) int); ok {
		r0 = rf(ctx, userId, userIdNew, coinId, userAccountTypeId, minpay)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int32, uuid.UUID, int32, int32, float64) error); ok {
		r1 = rf(ctx, userId, userIdNew, coinId, userAccountTypeId, minpay)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserAccountService_CreateUserAccount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUserAccount'
type MockUserAccountService_CreateUserAccount_Call struct {
	*mock.Call
}

// CreateUserAccount is a helper method to define mock.On call
//   - ctx context.Context
//   - userId int32
//   - userIdNew uuid.UUID
//   - coinId int32
//   - userAccountTypeId int32
//   - minpay float64
func (_e *MockUserAccountService_Expecter) CreateUserAccount(ctx interface{}, userId interface{}, userIdNew interface{}, coinId interface{}, userAccountTypeId interface{}, minpay interface{}) *MockUserAccountService_CreateUserAccount_Call {
	return &MockUserAccountService_CreateUserAccount_Call{Call: _e.mock.On("CreateUserAccount", ctx, userId, userIdNew, coinId, userAccountTypeId, minpay)}
}

func (_c *MockUserAccountService_CreateUserAccount_Call) Run(run func(ctx context.Context, userId int32, userIdNew uuid.UUID, coinId int32, userAccountTypeId int32, minpay float64)) *MockUserAccountService_CreateUserAccount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int32), args[2].(uuid.UUID), args[3].(int32), args[4].(int32), args[5].(float64))
	})
	return _c
}

func (_c *MockUserAccountService_CreateUserAccount_Call) Return(_a0 int, _a1 error) *MockUserAccountService_CreateUserAccount_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserAccountService_CreateUserAccount_Call) RunAndReturn(run func(context.Context, int32, uuid.UUID, int32, int32, float64) (int, error)) *MockUserAccountService_CreateUserAccount_Call {
	_c.Call.Return(run)
	return _c
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

// GetUserAccountIdByLegacyParams provides a mock function with given fields: ctx, userId, coinId, userAccountTypeId
func (_m *MockUserAccountService) GetUserAccountIdByLegacyParams(ctx context.Context, userId int32, coinId int32, userAccountTypeId int32) (int, error) {
	ret := _m.Called(ctx, userId, coinId, userAccountTypeId)

	if len(ret) == 0 {
		panic("no return value specified for GetUserAccountIdByLegacyParams")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int32, int32, int32) (int, error)); ok {
		return rf(ctx, userId, coinId, userAccountTypeId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int32, int32, int32) int); ok {
		r0 = rf(ctx, userId, coinId, userAccountTypeId)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int32, int32, int32) error); ok {
		r1 = rf(ctx, userId, coinId, userAccountTypeId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserAccountService_GetUserAccountIdByLegacyParams_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserAccountIdByLegacyParams'
type MockUserAccountService_GetUserAccountIdByLegacyParams_Call struct {
	*mock.Call
}

// GetUserAccountIdByLegacyParams is a helper method to define mock.On call
//   - ctx context.Context
//   - userId int32
//   - coinId int32
//   - userAccountTypeId int32
func (_e *MockUserAccountService_Expecter) GetUserAccountIdByLegacyParams(ctx interface{}, userId interface{}, coinId interface{}, userAccountTypeId interface{}) *MockUserAccountService_GetUserAccountIdByLegacyParams_Call {
	return &MockUserAccountService_GetUserAccountIdByLegacyParams_Call{Call: _e.mock.On("GetUserAccountIdByLegacyParams", ctx, userId, coinId, userAccountTypeId)}
}

func (_c *MockUserAccountService_GetUserAccountIdByLegacyParams_Call) Run(run func(ctx context.Context, userId int32, coinId int32, userAccountTypeId int32)) *MockUserAccountService_GetUserAccountIdByLegacyParams_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int32), args[2].(int32), args[3].(int32))
	})
	return _c
}

func (_c *MockUserAccountService_GetUserAccountIdByLegacyParams_Call) Return(_a0 int, _a1 error) *MockUserAccountService_GetUserAccountIdByLegacyParams_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserAccountService_GetUserAccountIdByLegacyParams_Call) RunAndReturn(run func(context.Context, int32, int32, int32) (int, error)) *MockUserAccountService_GetUserAccountIdByLegacyParams_Call {
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

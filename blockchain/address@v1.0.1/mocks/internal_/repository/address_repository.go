// Code generated by mockery v2.43.2. DO NOT EDIT.

package repository

import (
	context "context"

	model "code.emcdtech.com/emcd/blockchain/address/model"
	mock "github.com/stretchr/testify/mock"

	pgx "github.com/jackc/pgx/v5"

	repository "code.emcdtech.com/emcd/blockchain/address/internal/repository"

	transactor "code.emcdtech.com/emcd/sdk/pg"
)

// MockAddressRepository is an autogenerated mock type for the AddressRepository type
type MockAddressRepository struct {
	mock.Mock
}

type MockAddressRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAddressRepository) EXPECT() *MockAddressRepository_Expecter {
	return &MockAddressRepository_Expecter{mock: &_m.Mock}
}

// AddNewCommonAddress provides a mock function with given fields: ctx, address
func (_m *MockAddressRepository) AddNewCommonAddress(ctx context.Context, address *model.Address) error {
	ret := _m.Called(ctx, address)

	if len(ret) == 0 {
		panic("no return value specified for AddNewCommonAddress")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Address) error); ok {
		r0 = rf(ctx, address)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAddressRepository_AddNewCommonAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddNewCommonAddress'
type MockAddressRepository_AddNewCommonAddress_Call struct {
	*mock.Call
}

// AddNewCommonAddress is a helper method to define mock.On call
//   - ctx context.Context
//   - address *model.Address
func (_e *MockAddressRepository_Expecter) AddNewCommonAddress(ctx interface{}, address interface{}) *MockAddressRepository_AddNewCommonAddress_Call {
	return &MockAddressRepository_AddNewCommonAddress_Call{Call: _e.mock.On("AddNewCommonAddress", ctx, address)}
}

func (_c *MockAddressRepository_AddNewCommonAddress_Call) Run(run func(ctx context.Context, address *model.Address)) *MockAddressRepository_AddNewCommonAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Address))
	})
	return _c
}

func (_c *MockAddressRepository_AddNewCommonAddress_Call) Return(_a0 error) *MockAddressRepository_AddNewCommonAddress_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAddressRepository_AddNewCommonAddress_Call) RunAndReturn(run func(context.Context, *model.Address) error) *MockAddressRepository_AddNewCommonAddress_Call {
	_c.Call.Return(run)
	return _c
}

// AddNewDerivedAddress provides a mock function with given fields: ctx, address, masterKeyId, derivedFunc
func (_m *MockAddressRepository) AddNewDerivedAddress(ctx context.Context, address *model.Address, masterKeyId uint32, derivedFunc repository.DerivedFunc) error {
	ret := _m.Called(ctx, address, masterKeyId, derivedFunc)

	if len(ret) == 0 {
		panic("no return value specified for AddNewDerivedAddress")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Address, uint32, repository.DerivedFunc) error); ok {
		r0 = rf(ctx, address, masterKeyId, derivedFunc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAddressRepository_AddNewDerivedAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddNewDerivedAddress'
type MockAddressRepository_AddNewDerivedAddress_Call struct {
	*mock.Call
}

// AddNewDerivedAddress is a helper method to define mock.On call
//   - ctx context.Context
//   - address *model.Address
//   - masterKeyId uint32
//   - derivedFunc repository.DerivedFunc
func (_e *MockAddressRepository_Expecter) AddNewDerivedAddress(ctx interface{}, address interface{}, masterKeyId interface{}, derivedFunc interface{}) *MockAddressRepository_AddNewDerivedAddress_Call {
	return &MockAddressRepository_AddNewDerivedAddress_Call{Call: _e.mock.On("AddNewDerivedAddress", ctx, address, masterKeyId, derivedFunc)}
}

func (_c *MockAddressRepository_AddNewDerivedAddress_Call) Run(run func(ctx context.Context, address *model.Address, masterKeyId uint32, derivedFunc repository.DerivedFunc)) *MockAddressRepository_AddNewDerivedAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Address), args[2].(uint32), args[3].(repository.DerivedFunc))
	})
	return _c
}

func (_c *MockAddressRepository_AddNewDerivedAddress_Call) Return(_a0 error) *MockAddressRepository_AddNewDerivedAddress_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAddressRepository_AddNewDerivedAddress_Call) RunAndReturn(run func(context.Context, *model.Address, uint32, repository.DerivedFunc) error) *MockAddressRepository_AddNewDerivedAddress_Call {
	_c.Call.Return(run)
	return _c
}

// AddOldAddress provides a mock function with given fields: ctx, address
func (_m *MockAddressRepository) AddOldAddress(ctx context.Context, address *model.AddressOld) error {
	ret := _m.Called(ctx, address)

	if len(ret) == 0 {
		panic("no return value specified for AddOldAddress")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AddressOld) error); ok {
		r0 = rf(ctx, address)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAddressRepository_AddOldAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddOldAddress'
type MockAddressRepository_AddOldAddress_Call struct {
	*mock.Call
}

// AddOldAddress is a helper method to define mock.On call
//   - ctx context.Context
//   - address *model.AddressOld
func (_e *MockAddressRepository_Expecter) AddOldAddress(ctx interface{}, address interface{}) *MockAddressRepository_AddOldAddress_Call {
	return &MockAddressRepository_AddOldAddress_Call{Call: _e.mock.On("AddOldAddress", ctx, address)}
}

func (_c *MockAddressRepository_AddOldAddress_Call) Run(run func(ctx context.Context, address *model.AddressOld)) *MockAddressRepository_AddOldAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.AddressOld))
	})
	return _c
}

func (_c *MockAddressRepository_AddOldAddress_Call) Return(_a0 error) *MockAddressRepository_AddOldAddress_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAddressRepository_AddOldAddress_Call) RunAndReturn(run func(context.Context, *model.AddressOld) error) *MockAddressRepository_AddOldAddress_Call {
	_c.Call.Return(run)
	return _c
}

// AddOrUpdateDirtyAddress provides a mock function with given fields: ctx, address
func (_m *MockAddressRepository) AddOrUpdateDirtyAddress(ctx context.Context, address *model.AddressDirty) error {
	ret := _m.Called(ctx, address)

	if len(ret) == 0 {
		panic("no return value specified for AddOrUpdateDirtyAddress")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AddressDirty) error); ok {
		r0 = rf(ctx, address)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAddressRepository_AddOrUpdateDirtyAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddOrUpdateDirtyAddress'
type MockAddressRepository_AddOrUpdateDirtyAddress_Call struct {
	*mock.Call
}

// AddOrUpdateDirtyAddress is a helper method to define mock.On call
//   - ctx context.Context
//   - address *model.AddressDirty
func (_e *MockAddressRepository_Expecter) AddOrUpdateDirtyAddress(ctx interface{}, address interface{}) *MockAddressRepository_AddOrUpdateDirtyAddress_Call {
	return &MockAddressRepository_AddOrUpdateDirtyAddress_Call{Call: _e.mock.On("AddOrUpdateDirtyAddress", ctx, address)}
}

func (_c *MockAddressRepository_AddOrUpdateDirtyAddress_Call) Run(run func(ctx context.Context, address *model.AddressDirty)) *MockAddressRepository_AddOrUpdateDirtyAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.AddressDirty))
	})
	return _c
}

func (_c *MockAddressRepository_AddOrUpdateDirtyAddress_Call) Return(_a0 error) *MockAddressRepository_AddOrUpdateDirtyAddress_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAddressRepository_AddOrUpdateDirtyAddress_Call) RunAndReturn(run func(context.Context, *model.AddressDirty) error) *MockAddressRepository_AddOrUpdateDirtyAddress_Call {
	_c.Call.Return(run)
	return _c
}

// AddPersonalAddress provides a mock function with given fields: ctx, address
func (_m *MockAddressRepository) AddPersonalAddress(ctx context.Context, address *model.AddressPersonal) error {
	ret := _m.Called(ctx, address)

	if len(ret) == 0 {
		panic("no return value specified for AddPersonalAddress")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AddressPersonal) error); ok {
		r0 = rf(ctx, address)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAddressRepository_AddPersonalAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddPersonalAddress'
type MockAddressRepository_AddPersonalAddress_Call struct {
	*mock.Call
}

// AddPersonalAddress is a helper method to define mock.On call
//   - ctx context.Context
//   - address *model.AddressPersonal
func (_e *MockAddressRepository_Expecter) AddPersonalAddress(ctx interface{}, address interface{}) *MockAddressRepository_AddPersonalAddress_Call {
	return &MockAddressRepository_AddPersonalAddress_Call{Call: _e.mock.On("AddPersonalAddress", ctx, address)}
}

func (_c *MockAddressRepository_AddPersonalAddress_Call) Run(run func(ctx context.Context, address *model.AddressPersonal)) *MockAddressRepository_AddPersonalAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.AddressPersonal))
	})
	return _c
}

func (_c *MockAddressRepository_AddPersonalAddress_Call) Return(_a0 error) *MockAddressRepository_AddPersonalAddress_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAddressRepository_AddPersonalAddress_Call) RunAndReturn(run func(context.Context, *model.AddressPersonal) error) *MockAddressRepository_AddPersonalAddress_Call {
	_c.Call.Return(run)
	return _c
}

// GetDirtyAddresses provides a mock function with given fields: ctx, addressFilter
func (_m *MockAddressRepository) GetDirtyAddresses(ctx context.Context, addressFilter *model.AddressDirtyFilter) (model.AddressesDirty, error) {
	ret := _m.Called(ctx, addressFilter)

	if len(ret) == 0 {
		panic("no return value specified for GetDirtyAddresses")
	}

	var r0 model.AddressesDirty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AddressDirtyFilter) (model.AddressesDirty, error)); ok {
		return rf(ctx, addressFilter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.AddressDirtyFilter) model.AddressesDirty); ok {
		r0 = rf(ctx, addressFilter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.AddressesDirty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.AddressDirtyFilter) error); ok {
		r1 = rf(ctx, addressFilter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAddressRepository_GetDirtyAddresses_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDirtyAddresses'
type MockAddressRepository_GetDirtyAddresses_Call struct {
	*mock.Call
}

// GetDirtyAddresses is a helper method to define mock.On call
//   - ctx context.Context
//   - addressFilter *model.AddressDirtyFilter
func (_e *MockAddressRepository_Expecter) GetDirtyAddresses(ctx interface{}, addressFilter interface{}) *MockAddressRepository_GetDirtyAddresses_Call {
	return &MockAddressRepository_GetDirtyAddresses_Call{Call: _e.mock.On("GetDirtyAddresses", ctx, addressFilter)}
}

func (_c *MockAddressRepository_GetDirtyAddresses_Call) Run(run func(ctx context.Context, addressFilter *model.AddressDirtyFilter)) *MockAddressRepository_GetDirtyAddresses_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.AddressDirtyFilter))
	})
	return _c
}

func (_c *MockAddressRepository_GetDirtyAddresses_Call) Return(_a0 model.AddressesDirty, _a1 error) *MockAddressRepository_GetDirtyAddresses_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAddressRepository_GetDirtyAddresses_Call) RunAndReturn(run func(context.Context, *model.AddressDirtyFilter) (model.AddressesDirty, error)) *MockAddressRepository_GetDirtyAddresses_Call {
	_c.Call.Return(run)
	return _c
}

// GetNewAddresses provides a mock function with given fields: ctx, addressFilter
func (_m *MockAddressRepository) GetNewAddresses(ctx context.Context, addressFilter *model.AddressFilter) (*uint64, model.Addresses, error) {
	ret := _m.Called(ctx, addressFilter)

	if len(ret) == 0 {
		panic("no return value specified for GetNewAddresses")
	}

	var r0 *uint64
	var r1 model.Addresses
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AddressFilter) (*uint64, model.Addresses, error)); ok {
		return rf(ctx, addressFilter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.AddressFilter) *uint64); ok {
		r0 = rf(ctx, addressFilter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uint64)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.AddressFilter) model.Addresses); ok {
		r1 = rf(ctx, addressFilter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(model.Addresses)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, *model.AddressFilter) error); ok {
		r2 = rf(ctx, addressFilter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockAddressRepository_GetNewAddresses_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNewAddresses'
type MockAddressRepository_GetNewAddresses_Call struct {
	*mock.Call
}

// GetNewAddresses is a helper method to define mock.On call
//   - ctx context.Context
//   - addressFilter *model.AddressFilter
func (_e *MockAddressRepository_Expecter) GetNewAddresses(ctx interface{}, addressFilter interface{}) *MockAddressRepository_GetNewAddresses_Call {
	return &MockAddressRepository_GetNewAddresses_Call{Call: _e.mock.On("GetNewAddresses", ctx, addressFilter)}
}

func (_c *MockAddressRepository_GetNewAddresses_Call) Run(run func(ctx context.Context, addressFilter *model.AddressFilter)) *MockAddressRepository_GetNewAddresses_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.AddressFilter))
	})
	return _c
}

func (_c *MockAddressRepository_GetNewAddresses_Call) Return(_a0 *uint64, _a1 model.Addresses, _a2 error) *MockAddressRepository_GetNewAddresses_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockAddressRepository_GetNewAddresses_Call) RunAndReturn(run func(context.Context, *model.AddressFilter) (*uint64, model.Addresses, error)) *MockAddressRepository_GetNewAddresses_Call {
	_c.Call.Return(run)
	return _c
}

// GetOldAddresses provides a mock function with given fields: ctx, addressOldFilter
func (_m *MockAddressRepository) GetOldAddresses(ctx context.Context, addressOldFilter *model.AddressOldFilter) (*uint64, model.AddressesOld, error) {
	ret := _m.Called(ctx, addressOldFilter)

	if len(ret) == 0 {
		panic("no return value specified for GetOldAddresses")
	}

	var r0 *uint64
	var r1 model.AddressesOld
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AddressOldFilter) (*uint64, model.AddressesOld, error)); ok {
		return rf(ctx, addressOldFilter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.AddressOldFilter) *uint64); ok {
		r0 = rf(ctx, addressOldFilter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uint64)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.AddressOldFilter) model.AddressesOld); ok {
		r1 = rf(ctx, addressOldFilter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(model.AddressesOld)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, *model.AddressOldFilter) error); ok {
		r2 = rf(ctx, addressOldFilter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockAddressRepository_GetOldAddresses_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOldAddresses'
type MockAddressRepository_GetOldAddresses_Call struct {
	*mock.Call
}

// GetOldAddresses is a helper method to define mock.On call
//   - ctx context.Context
//   - addressOldFilter *model.AddressOldFilter
func (_e *MockAddressRepository_Expecter) GetOldAddresses(ctx interface{}, addressOldFilter interface{}) *MockAddressRepository_GetOldAddresses_Call {
	return &MockAddressRepository_GetOldAddresses_Call{Call: _e.mock.On("GetOldAddresses", ctx, addressOldFilter)}
}

func (_c *MockAddressRepository_GetOldAddresses_Call) Run(run func(ctx context.Context, addressOldFilter *model.AddressOldFilter)) *MockAddressRepository_GetOldAddresses_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.AddressOldFilter))
	})
	return _c
}

func (_c *MockAddressRepository_GetOldAddresses_Call) Return(_a0 *uint64, _a1 model.AddressesOld, _a2 error) *MockAddressRepository_GetOldAddresses_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockAddressRepository_GetOldAddresses_Call) RunAndReturn(run func(context.Context, *model.AddressOldFilter) (*uint64, model.AddressesOld, error)) *MockAddressRepository_GetOldAddresses_Call {
	_c.Call.Return(run)
	return _c
}

// GetPersonalAddresses provides a mock function with given fields: ctx, addressFilter
func (_m *MockAddressRepository) GetPersonalAddresses(ctx context.Context, addressFilter *model.AddressPersonalFilter) (*uint64, model.AddressesPersonal, error) {
	ret := _m.Called(ctx, addressFilter)

	if len(ret) == 0 {
		panic("no return value specified for GetPersonalAddresses")
	}

	var r0 *uint64
	var r1 model.AddressesPersonal
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AddressPersonalFilter) (*uint64, model.AddressesPersonal, error)); ok {
		return rf(ctx, addressFilter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.AddressPersonalFilter) *uint64); ok {
		r0 = rf(ctx, addressFilter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uint64)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.AddressPersonalFilter) model.AddressesPersonal); ok {
		r1 = rf(ctx, addressFilter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(model.AddressesPersonal)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, *model.AddressPersonalFilter) error); ok {
		r2 = rf(ctx, addressFilter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockAddressRepository_GetPersonalAddresses_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPersonalAddresses'
type MockAddressRepository_GetPersonalAddresses_Call struct {
	*mock.Call
}

// GetPersonalAddresses is a helper method to define mock.On call
//   - ctx context.Context
//   - addressFilter *model.AddressPersonalFilter
func (_e *MockAddressRepository_Expecter) GetPersonalAddresses(ctx interface{}, addressFilter interface{}) *MockAddressRepository_GetPersonalAddresses_Call {
	return &MockAddressRepository_GetPersonalAddresses_Call{Call: _e.mock.On("GetPersonalAddresses", ctx, addressFilter)}
}

func (_c *MockAddressRepository_GetPersonalAddresses_Call) Run(run func(ctx context.Context, addressFilter *model.AddressPersonalFilter)) *MockAddressRepository_GetPersonalAddresses_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.AddressPersonalFilter))
	})
	return _c
}

func (_c *MockAddressRepository_GetPersonalAddresses_Call) Return(_a0 *uint64, _a1 model.AddressesPersonal, _a2 error) *MockAddressRepository_GetPersonalAddresses_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockAddressRepository_GetPersonalAddresses_Call) RunAndReturn(run func(context.Context, *model.AddressPersonalFilter) (*uint64, model.AddressesPersonal, error)) *MockAddressRepository_GetPersonalAddresses_Call {
	_c.Call.Return(run)
	return _c
}

// Runner provides a mock function with given fields: ctx
func (_m *MockAddressRepository) Runner(ctx context.Context) transactor.PgxQueryRunner {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Runner")
	}

	var r0 transactor.PgxQueryRunner
	if rf, ok := ret.Get(0).(func(context.Context) transactor.PgxQueryRunner); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(transactor.PgxQueryRunner)
		}
	}

	return r0
}

// MockAddressRepository_Runner_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Runner'
type MockAddressRepository_Runner_Call struct {
	*mock.Call
}

// Runner is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockAddressRepository_Expecter) Runner(ctx interface{}) *MockAddressRepository_Runner_Call {
	return &MockAddressRepository_Runner_Call{Call: _e.mock.On("Runner", ctx)}
}

func (_c *MockAddressRepository_Runner_Call) Run(run func(ctx context.Context)) *MockAddressRepository_Runner_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockAddressRepository_Runner_Call) Return(_a0 transactor.PgxQueryRunner) *MockAddressRepository_Runner_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAddressRepository_Runner_Call) RunAndReturn(run func(context.Context) transactor.PgxQueryRunner) *MockAddressRepository_Runner_Call {
	_c.Call.Return(run)
	return _c
}

// UpdatePersonalAddress provides a mock function with given fields: ctx, address, addressPartial
func (_m *MockAddressRepository) UpdatePersonalAddress(ctx context.Context, address *model.AddressPersonal, addressPartial *model.AddressPersonalPartial) error {
	ret := _m.Called(ctx, address, addressPartial)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePersonalAddress")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AddressPersonal, *model.AddressPersonalPartial) error); ok {
		r0 = rf(ctx, address, addressPartial)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAddressRepository_UpdatePersonalAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdatePersonalAddress'
type MockAddressRepository_UpdatePersonalAddress_Call struct {
	*mock.Call
}

// UpdatePersonalAddress is a helper method to define mock.On call
//   - ctx context.Context
//   - address *model.AddressPersonal
//   - addressPartial *model.AddressPersonalPartial
func (_e *MockAddressRepository_Expecter) UpdatePersonalAddress(ctx interface{}, address interface{}, addressPartial interface{}) *MockAddressRepository_UpdatePersonalAddress_Call {
	return &MockAddressRepository_UpdatePersonalAddress_Call{Call: _e.mock.On("UpdatePersonalAddress", ctx, address, addressPartial)}
}

func (_c *MockAddressRepository_UpdatePersonalAddress_Call) Run(run func(ctx context.Context, address *model.AddressPersonal, addressPartial *model.AddressPersonalPartial)) *MockAddressRepository_UpdatePersonalAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.AddressPersonal), args[2].(*model.AddressPersonalPartial))
	})
	return _c
}

func (_c *MockAddressRepository_UpdatePersonalAddress_Call) Return(_a0 error) *MockAddressRepository_UpdatePersonalAddress_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAddressRepository_UpdatePersonalAddress_Call) RunAndReturn(run func(context.Context, *model.AddressPersonal, *model.AddressPersonalPartial) error) *MockAddressRepository_UpdatePersonalAddress_Call {
	_c.Call.Return(run)
	return _c
}

// WithinTransaction provides a mock function with given fields: ctx, txFn
func (_m *MockAddressRepository) WithinTransaction(ctx context.Context, txFn func(context.Context) error) error {
	ret := _m.Called(ctx, txFn)

	if len(ret) == 0 {
		panic("no return value specified for WithinTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context) error) error); ok {
		r0 = rf(ctx, txFn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAddressRepository_WithinTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithinTransaction'
type MockAddressRepository_WithinTransaction_Call struct {
	*mock.Call
}

// WithinTransaction is a helper method to define mock.On call
//   - ctx context.Context
//   - txFn func(context.Context) error
func (_e *MockAddressRepository_Expecter) WithinTransaction(ctx interface{}, txFn interface{}) *MockAddressRepository_WithinTransaction_Call {
	return &MockAddressRepository_WithinTransaction_Call{Call: _e.mock.On("WithinTransaction", ctx, txFn)}
}

func (_c *MockAddressRepository_WithinTransaction_Call) Run(run func(ctx context.Context, txFn func(context.Context) error)) *MockAddressRepository_WithinTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(context.Context) error))
	})
	return _c
}

func (_c *MockAddressRepository_WithinTransaction_Call) Return(_a0 error) *MockAddressRepository_WithinTransaction_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAddressRepository_WithinTransaction_Call) RunAndReturn(run func(context.Context, func(context.Context) error) error) *MockAddressRepository_WithinTransaction_Call {
	_c.Call.Return(run)
	return _c
}

// WithinTransactionWithOptions provides a mock function with given fields: ctx, txFn, opts
func (_m *MockAddressRepository) WithinTransactionWithOptions(ctx context.Context, txFn func(context.Context) error, opts pgx.TxOptions) error {
	ret := _m.Called(ctx, txFn, opts)

	if len(ret) == 0 {
		panic("no return value specified for WithinTransactionWithOptions")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context) error, pgx.TxOptions) error); ok {
		r0 = rf(ctx, txFn, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAddressRepository_WithinTransactionWithOptions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithinTransactionWithOptions'
type MockAddressRepository_WithinTransactionWithOptions_Call struct {
	*mock.Call
}

// WithinTransactionWithOptions is a helper method to define mock.On call
//   - ctx context.Context
//   - txFn func(context.Context) error
//   - opts pgx.TxOptions
func (_e *MockAddressRepository_Expecter) WithinTransactionWithOptions(ctx interface{}, txFn interface{}, opts interface{}) *MockAddressRepository_WithinTransactionWithOptions_Call {
	return &MockAddressRepository_WithinTransactionWithOptions_Call{Call: _e.mock.On("WithinTransactionWithOptions", ctx, txFn, opts)}
}

func (_c *MockAddressRepository_WithinTransactionWithOptions_Call) Run(run func(ctx context.Context, txFn func(context.Context) error, opts pgx.TxOptions)) *MockAddressRepository_WithinTransactionWithOptions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(context.Context) error), args[2].(pgx.TxOptions))
	})
	return _c
}

func (_c *MockAddressRepository_WithinTransactionWithOptions_Call) Return(_a0 error) *MockAddressRepository_WithinTransactionWithOptions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAddressRepository_WithinTransactionWithOptions_Call) RunAndReturn(run func(context.Context, func(context.Context) error, pgx.TxOptions) error) *MockAddressRepository_WithinTransactionWithOptions_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAddressRepository creates a new instance of MockAddressRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAddressRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAddressRepository {
	mock := &MockAddressRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
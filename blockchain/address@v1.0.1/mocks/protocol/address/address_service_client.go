// Code generated by mockery v2.43.2. DO NOT EDIT.

package address

import (
	context "context"

	address "code.emcdtech.com/emcd/blockchain/address/protocol/address"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// MockAddressServiceClient is an autogenerated mock type for the AddressServiceClient type
type MockAddressServiceClient struct {
	mock.Mock
}

type MockAddressServiceClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAddressServiceClient) EXPECT() *MockAddressServiceClient_Expecter {
	return &MockAddressServiceClient_Expecter{mock: &_m.Mock}
}

// AddOrUpdatePersonalAddress provides a mock function with given fields: ctx, in, opts
func (_m *MockAddressServiceClient) AddOrUpdatePersonalAddress(ctx context.Context, in *address.CreatePersonalAddressRequest, opts ...grpc.CallOption) (*address.PersonalAddressResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for AddOrUpdatePersonalAddress")
	}

	var r0 *address.PersonalAddressResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *address.CreatePersonalAddressRequest, ...grpc.CallOption) (*address.PersonalAddressResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *address.CreatePersonalAddressRequest, ...grpc.CallOption) *address.PersonalAddressResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*address.PersonalAddressResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *address.CreatePersonalAddressRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAddressServiceClient_AddOrUpdatePersonalAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddOrUpdatePersonalAddress'
type MockAddressServiceClient_AddOrUpdatePersonalAddress_Call struct {
	*mock.Call
}

// AddOrUpdatePersonalAddress is a helper method to define mock.On call
//   - ctx context.Context
//   - in *address.CreatePersonalAddressRequest
//   - opts ...grpc.CallOption
func (_e *MockAddressServiceClient_Expecter) AddOrUpdatePersonalAddress(ctx interface{}, in interface{}, opts ...interface{}) *MockAddressServiceClient_AddOrUpdatePersonalAddress_Call {
	return &MockAddressServiceClient_AddOrUpdatePersonalAddress_Call{Call: _e.mock.On("AddOrUpdatePersonalAddress",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAddressServiceClient_AddOrUpdatePersonalAddress_Call) Run(run func(ctx context.Context, in *address.CreatePersonalAddressRequest, opts ...grpc.CallOption)) *MockAddressServiceClient_AddOrUpdatePersonalAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*address.CreatePersonalAddressRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockAddressServiceClient_AddOrUpdatePersonalAddress_Call) Return(_a0 *address.PersonalAddressResponse, _a1 error) *MockAddressServiceClient_AddOrUpdatePersonalAddress_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAddressServiceClient_AddOrUpdatePersonalAddress_Call) RunAndReturn(run func(context.Context, *address.CreatePersonalAddressRequest, ...grpc.CallOption) (*address.PersonalAddressResponse, error)) *MockAddressServiceClient_AddOrUpdatePersonalAddress_Call {
	_c.Call.Return(run)
	return _c
}

// CreateOrUpdateDirtyAddress provides a mock function with given fields: ctx, in, opts
func (_m *MockAddressServiceClient) CreateOrUpdateDirtyAddress(ctx context.Context, in *address.DirtyAddressForm, opts ...grpc.CallOption) (*address.DirtyAddressForm, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrUpdateDirtyAddress")
	}

	var r0 *address.DirtyAddressForm
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *address.DirtyAddressForm, ...grpc.CallOption) (*address.DirtyAddressForm, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *address.DirtyAddressForm, ...grpc.CallOption) *address.DirtyAddressForm); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*address.DirtyAddressForm)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *address.DirtyAddressForm, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAddressServiceClient_CreateOrUpdateDirtyAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateOrUpdateDirtyAddress'
type MockAddressServiceClient_CreateOrUpdateDirtyAddress_Call struct {
	*mock.Call
}

// CreateOrUpdateDirtyAddress is a helper method to define mock.On call
//   - ctx context.Context
//   - in *address.DirtyAddressForm
//   - opts ...grpc.CallOption
func (_e *MockAddressServiceClient_Expecter) CreateOrUpdateDirtyAddress(ctx interface{}, in interface{}, opts ...interface{}) *MockAddressServiceClient_CreateOrUpdateDirtyAddress_Call {
	return &MockAddressServiceClient_CreateOrUpdateDirtyAddress_Call{Call: _e.mock.On("CreateOrUpdateDirtyAddress",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAddressServiceClient_CreateOrUpdateDirtyAddress_Call) Run(run func(ctx context.Context, in *address.DirtyAddressForm, opts ...grpc.CallOption)) *MockAddressServiceClient_CreateOrUpdateDirtyAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*address.DirtyAddressForm), variadicArgs...)
	})
	return _c
}

func (_c *MockAddressServiceClient_CreateOrUpdateDirtyAddress_Call) Return(_a0 *address.DirtyAddressForm, _a1 error) *MockAddressServiceClient_CreateOrUpdateDirtyAddress_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAddressServiceClient_CreateOrUpdateDirtyAddress_Call) RunAndReturn(run func(context.Context, *address.DirtyAddressForm, ...grpc.CallOption) (*address.DirtyAddressForm, error)) *MockAddressServiceClient_CreateOrUpdateDirtyAddress_Call {
	_c.Call.Return(run)
	return _c
}

// CreateProcessingAddress provides a mock function with given fields: ctx, in, opts
func (_m *MockAddressServiceClient) CreateProcessingAddress(ctx context.Context, in *address.CreateProcessingAddressRequest, opts ...grpc.CallOption) (*address.AddressResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CreateProcessingAddress")
	}

	var r0 *address.AddressResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *address.CreateProcessingAddressRequest, ...grpc.CallOption) (*address.AddressResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *address.CreateProcessingAddressRequest, ...grpc.CallOption) *address.AddressResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*address.AddressResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *address.CreateProcessingAddressRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAddressServiceClient_CreateProcessingAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateProcessingAddress'
type MockAddressServiceClient_CreateProcessingAddress_Call struct {
	*mock.Call
}

// CreateProcessingAddress is a helper method to define mock.On call
//   - ctx context.Context
//   - in *address.CreateProcessingAddressRequest
//   - opts ...grpc.CallOption
func (_e *MockAddressServiceClient_Expecter) CreateProcessingAddress(ctx interface{}, in interface{}, opts ...interface{}) *MockAddressServiceClient_CreateProcessingAddress_Call {
	return &MockAddressServiceClient_CreateProcessingAddress_Call{Call: _e.mock.On("CreateProcessingAddress",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAddressServiceClient_CreateProcessingAddress_Call) Run(run func(ctx context.Context, in *address.CreateProcessingAddressRequest, opts ...grpc.CallOption)) *MockAddressServiceClient_CreateProcessingAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*address.CreateProcessingAddressRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockAddressServiceClient_CreateProcessingAddress_Call) Return(_a0 *address.AddressResponse, _a1 error) *MockAddressServiceClient_CreateProcessingAddress_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAddressServiceClient_CreateProcessingAddress_Call) RunAndReturn(run func(context.Context, *address.CreateProcessingAddressRequest, ...grpc.CallOption) (*address.AddressResponse, error)) *MockAddressServiceClient_CreateProcessingAddress_Call {
	_c.Call.Return(run)
	return _c
}

// DeletePersonalAddress provides a mock function with given fields: ctx, in, opts
func (_m *MockAddressServiceClient) DeletePersonalAddress(ctx context.Context, in *address.DeletePersonalAddressRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeletePersonalAddress")
	}

	var r0 *emptypb.Empty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *address.DeletePersonalAddressRequest, ...grpc.CallOption) (*emptypb.Empty, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *address.DeletePersonalAddressRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *address.DeletePersonalAddressRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAddressServiceClient_DeletePersonalAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeletePersonalAddress'
type MockAddressServiceClient_DeletePersonalAddress_Call struct {
	*mock.Call
}

// DeletePersonalAddress is a helper method to define mock.On call
//   - ctx context.Context
//   - in *address.DeletePersonalAddressRequest
//   - opts ...grpc.CallOption
func (_e *MockAddressServiceClient_Expecter) DeletePersonalAddress(ctx interface{}, in interface{}, opts ...interface{}) *MockAddressServiceClient_DeletePersonalAddress_Call {
	return &MockAddressServiceClient_DeletePersonalAddress_Call{Call: _e.mock.On("DeletePersonalAddress",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAddressServiceClient_DeletePersonalAddress_Call) Run(run func(ctx context.Context, in *address.DeletePersonalAddressRequest, opts ...grpc.CallOption)) *MockAddressServiceClient_DeletePersonalAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*address.DeletePersonalAddressRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockAddressServiceClient_DeletePersonalAddress_Call) Return(_a0 *emptypb.Empty, _a1 error) *MockAddressServiceClient_DeletePersonalAddress_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAddressServiceClient_DeletePersonalAddress_Call) RunAndReturn(run func(context.Context, *address.DeletePersonalAddressRequest, ...grpc.CallOption) (*emptypb.Empty, error)) *MockAddressServiceClient_DeletePersonalAddress_Call {
	_c.Call.Return(run)
	return _c
}

// GetAddressByStr provides a mock function with given fields: ctx, in, opts
func (_m *MockAddressServiceClient) GetAddressByStr(ctx context.Context, in *address.AddressStrId, opts ...grpc.CallOption) (*address.AddressResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetAddressByStr")
	}

	var r0 *address.AddressResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *address.AddressStrId, ...grpc.CallOption) (*address.AddressResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *address.AddressStrId, ...grpc.CallOption) *address.AddressResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*address.AddressResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *address.AddressStrId, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAddressServiceClient_GetAddressByStr_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAddressByStr'
type MockAddressServiceClient_GetAddressByStr_Call struct {
	*mock.Call
}

// GetAddressByStr is a helper method to define mock.On call
//   - ctx context.Context
//   - in *address.AddressStrId
//   - opts ...grpc.CallOption
func (_e *MockAddressServiceClient_Expecter) GetAddressByStr(ctx interface{}, in interface{}, opts ...interface{}) *MockAddressServiceClient_GetAddressByStr_Call {
	return &MockAddressServiceClient_GetAddressByStr_Call{Call: _e.mock.On("GetAddressByStr",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAddressServiceClient_GetAddressByStr_Call) Run(run func(ctx context.Context, in *address.AddressStrId, opts ...grpc.CallOption)) *MockAddressServiceClient_GetAddressByStr_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*address.AddressStrId), variadicArgs...)
	})
	return _c
}

func (_c *MockAddressServiceClient_GetAddressByStr_Call) Return(_a0 *address.AddressResponse, _a1 error) *MockAddressServiceClient_GetAddressByStr_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAddressServiceClient_GetAddressByStr_Call) RunAndReturn(run func(context.Context, *address.AddressStrId, ...grpc.CallOption) (*address.AddressResponse, error)) *MockAddressServiceClient_GetAddressByStr_Call {
	_c.Call.Return(run)
	return _c
}

// GetAddressByUuid provides a mock function with given fields: ctx, in, opts
func (_m *MockAddressServiceClient) GetAddressByUuid(ctx context.Context, in *address.AddressUuid, opts ...grpc.CallOption) (*address.AddressResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetAddressByUuid")
	}

	var r0 *address.AddressResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *address.AddressUuid, ...grpc.CallOption) (*address.AddressResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *address.AddressUuid, ...grpc.CallOption) *address.AddressResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*address.AddressResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *address.AddressUuid, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAddressServiceClient_GetAddressByUuid_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAddressByUuid'
type MockAddressServiceClient_GetAddressByUuid_Call struct {
	*mock.Call
}

// GetAddressByUuid is a helper method to define mock.On call
//   - ctx context.Context
//   - in *address.AddressUuid
//   - opts ...grpc.CallOption
func (_e *MockAddressServiceClient_Expecter) GetAddressByUuid(ctx interface{}, in interface{}, opts ...interface{}) *MockAddressServiceClient_GetAddressByUuid_Call {
	return &MockAddressServiceClient_GetAddressByUuid_Call{Call: _e.mock.On("GetAddressByUuid",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAddressServiceClient_GetAddressByUuid_Call) Run(run func(ctx context.Context, in *address.AddressUuid, opts ...grpc.CallOption)) *MockAddressServiceClient_GetAddressByUuid_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*address.AddressUuid), variadicArgs...)
	})
	return _c
}

func (_c *MockAddressServiceClient_GetAddressByUuid_Call) Return(_a0 *address.AddressResponse, _a1 error) *MockAddressServiceClient_GetAddressByUuid_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAddressServiceClient_GetAddressByUuid_Call) RunAndReturn(run func(context.Context, *address.AddressUuid, ...grpc.CallOption) (*address.AddressResponse, error)) *MockAddressServiceClient_GetAddressByUuid_Call {
	_c.Call.Return(run)
	return _c
}

// GetAddressesByUserUuid provides a mock function with given fields: ctx, in, opts
func (_m *MockAddressServiceClient) GetAddressesByUserUuid(ctx context.Context, in *address.UserUuid, opts ...grpc.CallOption) (*address.AddressMultiResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetAddressesByUserUuid")
	}

	var r0 *address.AddressMultiResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *address.UserUuid, ...grpc.CallOption) (*address.AddressMultiResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *address.UserUuid, ...grpc.CallOption) *address.AddressMultiResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*address.AddressMultiResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *address.UserUuid, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAddressServiceClient_GetAddressesByUserUuid_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAddressesByUserUuid'
type MockAddressServiceClient_GetAddressesByUserUuid_Call struct {
	*mock.Call
}

// GetAddressesByUserUuid is a helper method to define mock.On call
//   - ctx context.Context
//   - in *address.UserUuid
//   - opts ...grpc.CallOption
func (_e *MockAddressServiceClient_Expecter) GetAddressesByUserUuid(ctx interface{}, in interface{}, opts ...interface{}) *MockAddressServiceClient_GetAddressesByUserUuid_Call {
	return &MockAddressServiceClient_GetAddressesByUserUuid_Call{Call: _e.mock.On("GetAddressesByUserUuid",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAddressServiceClient_GetAddressesByUserUuid_Call) Run(run func(ctx context.Context, in *address.UserUuid, opts ...grpc.CallOption)) *MockAddressServiceClient_GetAddressesByUserUuid_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*address.UserUuid), variadicArgs...)
	})
	return _c
}

func (_c *MockAddressServiceClient_GetAddressesByUserUuid_Call) Return(_a0 *address.AddressMultiResponse, _a1 error) *MockAddressServiceClient_GetAddressesByUserUuid_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAddressServiceClient_GetAddressesByUserUuid_Call) RunAndReturn(run func(context.Context, *address.UserUuid, ...grpc.CallOption) (*address.AddressMultiResponse, error)) *MockAddressServiceClient_GetAddressesByUserUuid_Call {
	_c.Call.Return(run)
	return _c
}

// GetAddressesNewByFilter provides a mock function with given fields: ctx, in, opts
func (_m *MockAddressServiceClient) GetAddressesNewByFilter(ctx context.Context, in *address.AddressNewFilter, opts ...grpc.CallOption) (*address.AddressMultiResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetAddressesNewByFilter")
	}

	var r0 *address.AddressMultiResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *address.AddressNewFilter, ...grpc.CallOption) (*address.AddressMultiResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *address.AddressNewFilter, ...grpc.CallOption) *address.AddressMultiResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*address.AddressMultiResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *address.AddressNewFilter, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAddressServiceClient_GetAddressesNewByFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAddressesNewByFilter'
type MockAddressServiceClient_GetAddressesNewByFilter_Call struct {
	*mock.Call
}

// GetAddressesNewByFilter is a helper method to define mock.On call
//   - ctx context.Context
//   - in *address.AddressNewFilter
//   - opts ...grpc.CallOption
func (_e *MockAddressServiceClient_Expecter) GetAddressesNewByFilter(ctx interface{}, in interface{}, opts ...interface{}) *MockAddressServiceClient_GetAddressesNewByFilter_Call {
	return &MockAddressServiceClient_GetAddressesNewByFilter_Call{Call: _e.mock.On("GetAddressesNewByFilter",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAddressServiceClient_GetAddressesNewByFilter_Call) Run(run func(ctx context.Context, in *address.AddressNewFilter, opts ...grpc.CallOption)) *MockAddressServiceClient_GetAddressesNewByFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*address.AddressNewFilter), variadicArgs...)
	})
	return _c
}

func (_c *MockAddressServiceClient_GetAddressesNewByFilter_Call) Return(_a0 *address.AddressMultiResponse, _a1 error) *MockAddressServiceClient_GetAddressesNewByFilter_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAddressServiceClient_GetAddressesNewByFilter_Call) RunAndReturn(run func(context.Context, *address.AddressNewFilter, ...grpc.CallOption) (*address.AddressMultiResponse, error)) *MockAddressServiceClient_GetAddressesNewByFilter_Call {
	_c.Call.Return(run)
	return _c
}

// GetAddressesOldByFilter provides a mock function with given fields: ctx, in, opts
func (_m *MockAddressServiceClient) GetAddressesOldByFilter(ctx context.Context, in *address.AddressOldFilter, opts ...grpc.CallOption) (*address.AddressMultiResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetAddressesOldByFilter")
	}

	var r0 *address.AddressMultiResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *address.AddressOldFilter, ...grpc.CallOption) (*address.AddressMultiResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *address.AddressOldFilter, ...grpc.CallOption) *address.AddressMultiResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*address.AddressMultiResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *address.AddressOldFilter, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAddressServiceClient_GetAddressesOldByFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAddressesOldByFilter'
type MockAddressServiceClient_GetAddressesOldByFilter_Call struct {
	*mock.Call
}

// GetAddressesOldByFilter is a helper method to define mock.On call
//   - ctx context.Context
//   - in *address.AddressOldFilter
//   - opts ...grpc.CallOption
func (_e *MockAddressServiceClient_Expecter) GetAddressesOldByFilter(ctx interface{}, in interface{}, opts ...interface{}) *MockAddressServiceClient_GetAddressesOldByFilter_Call {
	return &MockAddressServiceClient_GetAddressesOldByFilter_Call{Call: _e.mock.On("GetAddressesOldByFilter",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAddressServiceClient_GetAddressesOldByFilter_Call) Run(run func(ctx context.Context, in *address.AddressOldFilter, opts ...grpc.CallOption)) *MockAddressServiceClient_GetAddressesOldByFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*address.AddressOldFilter), variadicArgs...)
	})
	return _c
}

func (_c *MockAddressServiceClient_GetAddressesOldByFilter_Call) Return(_a0 *address.AddressMultiResponse, _a1 error) *MockAddressServiceClient_GetAddressesOldByFilter_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAddressServiceClient_GetAddressesOldByFilter_Call) RunAndReturn(run func(context.Context, *address.AddressOldFilter, ...grpc.CallOption) (*address.AddressMultiResponse, error)) *MockAddressServiceClient_GetAddressesOldByFilter_Call {
	_c.Call.Return(run)
	return _c
}

// GetDirtyAddressesByFilter provides a mock function with given fields: ctx, in, opts
func (_m *MockAddressServiceClient) GetDirtyAddressesByFilter(ctx context.Context, in *address.DirtyAddressFilter, opts ...grpc.CallOption) (*address.DirtyAddressMultiForm, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetDirtyAddressesByFilter")
	}

	var r0 *address.DirtyAddressMultiForm
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *address.DirtyAddressFilter, ...grpc.CallOption) (*address.DirtyAddressMultiForm, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *address.DirtyAddressFilter, ...grpc.CallOption) *address.DirtyAddressMultiForm); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*address.DirtyAddressMultiForm)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *address.DirtyAddressFilter, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAddressServiceClient_GetDirtyAddressesByFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDirtyAddressesByFilter'
type MockAddressServiceClient_GetDirtyAddressesByFilter_Call struct {
	*mock.Call
}

// GetDirtyAddressesByFilter is a helper method to define mock.On call
//   - ctx context.Context
//   - in *address.DirtyAddressFilter
//   - opts ...grpc.CallOption
func (_e *MockAddressServiceClient_Expecter) GetDirtyAddressesByFilter(ctx interface{}, in interface{}, opts ...interface{}) *MockAddressServiceClient_GetDirtyAddressesByFilter_Call {
	return &MockAddressServiceClient_GetDirtyAddressesByFilter_Call{Call: _e.mock.On("GetDirtyAddressesByFilter",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAddressServiceClient_GetDirtyAddressesByFilter_Call) Run(run func(ctx context.Context, in *address.DirtyAddressFilter, opts ...grpc.CallOption)) *MockAddressServiceClient_GetDirtyAddressesByFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*address.DirtyAddressFilter), variadicArgs...)
	})
	return _c
}

func (_c *MockAddressServiceClient_GetDirtyAddressesByFilter_Call) Return(_a0 *address.DirtyAddressMultiForm, _a1 error) *MockAddressServiceClient_GetDirtyAddressesByFilter_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAddressServiceClient_GetDirtyAddressesByFilter_Call) RunAndReturn(run func(context.Context, *address.DirtyAddressFilter, ...grpc.CallOption) (*address.DirtyAddressMultiForm, error)) *MockAddressServiceClient_GetDirtyAddressesByFilter_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrCreateAddress provides a mock function with given fields: ctx, in, opts
func (_m *MockAddressServiceClient) GetOrCreateAddress(ctx context.Context, in *address.CreateAddressRequest, opts ...grpc.CallOption) (*address.AddressResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetOrCreateAddress")
	}

	var r0 *address.AddressResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *address.CreateAddressRequest, ...grpc.CallOption) (*address.AddressResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *address.CreateAddressRequest, ...grpc.CallOption) *address.AddressResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*address.AddressResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *address.CreateAddressRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAddressServiceClient_GetOrCreateAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrCreateAddress'
type MockAddressServiceClient_GetOrCreateAddress_Call struct {
	*mock.Call
}

// GetOrCreateAddress is a helper method to define mock.On call
//   - ctx context.Context
//   - in *address.CreateAddressRequest
//   - opts ...grpc.CallOption
func (_e *MockAddressServiceClient_Expecter) GetOrCreateAddress(ctx interface{}, in interface{}, opts ...interface{}) *MockAddressServiceClient_GetOrCreateAddress_Call {
	return &MockAddressServiceClient_GetOrCreateAddress_Call{Call: _e.mock.On("GetOrCreateAddress",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAddressServiceClient_GetOrCreateAddress_Call) Run(run func(ctx context.Context, in *address.CreateAddressRequest, opts ...grpc.CallOption)) *MockAddressServiceClient_GetOrCreateAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*address.CreateAddressRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockAddressServiceClient_GetOrCreateAddress_Call) Return(_a0 *address.AddressResponse, _a1 error) *MockAddressServiceClient_GetOrCreateAddress_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAddressServiceClient_GetOrCreateAddress_Call) RunAndReturn(run func(context.Context, *address.CreateAddressRequest, ...grpc.CallOption) (*address.AddressResponse, error)) *MockAddressServiceClient_GetOrCreateAddress_Call {
	_c.Call.Return(run)
	return _c
}

// GetPersonalAddressesByFilter provides a mock function with given fields: ctx, in, opts
func (_m *MockAddressServiceClient) GetPersonalAddressesByFilter(ctx context.Context, in *address.AddressPersonalFilter, opts ...grpc.CallOption) (*address.PersonalAddressMultiResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetPersonalAddressesByFilter")
	}

	var r0 *address.PersonalAddressMultiResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *address.AddressPersonalFilter, ...grpc.CallOption) (*address.PersonalAddressMultiResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *address.AddressPersonalFilter, ...grpc.CallOption) *address.PersonalAddressMultiResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*address.PersonalAddressMultiResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *address.AddressPersonalFilter, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAddressServiceClient_GetPersonalAddressesByFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPersonalAddressesByFilter'
type MockAddressServiceClient_GetPersonalAddressesByFilter_Call struct {
	*mock.Call
}

// GetPersonalAddressesByFilter is a helper method to define mock.On call
//   - ctx context.Context
//   - in *address.AddressPersonalFilter
//   - opts ...grpc.CallOption
func (_e *MockAddressServiceClient_Expecter) GetPersonalAddressesByFilter(ctx interface{}, in interface{}, opts ...interface{}) *MockAddressServiceClient_GetPersonalAddressesByFilter_Call {
	return &MockAddressServiceClient_GetPersonalAddressesByFilter_Call{Call: _e.mock.On("GetPersonalAddressesByFilter",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAddressServiceClient_GetPersonalAddressesByFilter_Call) Run(run func(ctx context.Context, in *address.AddressPersonalFilter, opts ...grpc.CallOption)) *MockAddressServiceClient_GetPersonalAddressesByFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*address.AddressPersonalFilter), variadicArgs...)
	})
	return _c
}

func (_c *MockAddressServiceClient_GetPersonalAddressesByFilter_Call) Return(_a0 *address.PersonalAddressMultiResponse, _a1 error) *MockAddressServiceClient_GetPersonalAddressesByFilter_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAddressServiceClient_GetPersonalAddressesByFilter_Call) RunAndReturn(run func(context.Context, *address.AddressPersonalFilter, ...grpc.CallOption) (*address.PersonalAddressMultiResponse, error)) *MockAddressServiceClient_GetPersonalAddressesByFilter_Call {
	_c.Call.Return(run)
	return _c
}

// GetPersonalAddressesByUserUuid provides a mock function with given fields: ctx, in, opts
func (_m *MockAddressServiceClient) GetPersonalAddressesByUserUuid(ctx context.Context, in *address.UserUuid, opts ...grpc.CallOption) (*address.PersonalAddressMultiResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetPersonalAddressesByUserUuid")
	}

	var r0 *address.PersonalAddressMultiResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *address.UserUuid, ...grpc.CallOption) (*address.PersonalAddressMultiResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *address.UserUuid, ...grpc.CallOption) *address.PersonalAddressMultiResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*address.PersonalAddressMultiResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *address.UserUuid, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAddressServiceClient_GetPersonalAddressesByUserUuid_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPersonalAddressesByUserUuid'
type MockAddressServiceClient_GetPersonalAddressesByUserUuid_Call struct {
	*mock.Call
}

// GetPersonalAddressesByUserUuid is a helper method to define mock.On call
//   - ctx context.Context
//   - in *address.UserUuid
//   - opts ...grpc.CallOption
func (_e *MockAddressServiceClient_Expecter) GetPersonalAddressesByUserUuid(ctx interface{}, in interface{}, opts ...interface{}) *MockAddressServiceClient_GetPersonalAddressesByUserUuid_Call {
	return &MockAddressServiceClient_GetPersonalAddressesByUserUuid_Call{Call: _e.mock.On("GetPersonalAddressesByUserUuid",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAddressServiceClient_GetPersonalAddressesByUserUuid_Call) Run(run func(ctx context.Context, in *address.UserUuid, opts ...grpc.CallOption)) *MockAddressServiceClient_GetPersonalAddressesByUserUuid_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*address.UserUuid), variadicArgs...)
	})
	return _c
}

func (_c *MockAddressServiceClient_GetPersonalAddressesByUserUuid_Call) Return(_a0 *address.PersonalAddressMultiResponse, _a1 error) *MockAddressServiceClient_GetPersonalAddressesByUserUuid_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAddressServiceClient_GetPersonalAddressesByUserUuid_Call) RunAndReturn(run func(context.Context, *address.UserUuid, ...grpc.CallOption) (*address.PersonalAddressMultiResponse, error)) *MockAddressServiceClient_GetPersonalAddressesByUserUuid_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAddressServiceClient creates a new instance of MockAddressServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAddressServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAddressServiceClient {
	mock := &MockAddressServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
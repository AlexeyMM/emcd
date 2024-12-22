// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	whitelabel "code.emcdtech.com/emcd/service/whitelabel/protocol/whitelabel"
)

// MockWhiteLabel is an autogenerated mock type for the WhiteLabel type
type MockWhiteLabel struct {
	mock.Mock
}

// GetCoins provides a mock function with given fields: ctx, in, opts
func (_m *MockWhiteLabel) GetCoins(ctx context.Context, in *whitelabel.GetCoinsRequest, opts ...grpc.CallOption) (*whitelabel.GetCoinsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetCoins")
	}

	var r0 *whitelabel.GetCoinsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *whitelabel.GetCoinsRequest, ...grpc.CallOption) (*whitelabel.GetCoinsResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *whitelabel.GetCoinsRequest, ...grpc.CallOption) *whitelabel.GetCoinsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*whitelabel.GetCoinsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *whitelabel.GetCoinsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockWhiteLabel creates a new instance of MockWhiteLabel. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockWhiteLabel(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockWhiteLabel {
	mock := &MockWhiteLabel{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

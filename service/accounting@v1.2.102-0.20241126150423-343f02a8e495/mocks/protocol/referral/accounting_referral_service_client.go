// Code generated by mockery v2.43.2. DO NOT EDIT.

package referral

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	referral "code.emcdtech.com/emcd/service/accounting/protocol/referral"
)

// MockAccountingReferralServiceClient is an autogenerated mock type for the AccountingReferralServiceClient type
type MockAccountingReferralServiceClient struct {
	mock.Mock
}

type MockAccountingReferralServiceClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAccountingReferralServiceClient) EXPECT() *MockAccountingReferralServiceClient_Expecter {
	return &MockAccountingReferralServiceClient_Expecter{mock: &_m.Mock}
}

// GetReferralsStatistic provides a mock function with given fields: ctx, in, opts
func (_m *MockAccountingReferralServiceClient) GetReferralsStatistic(ctx context.Context, in *referral.GetReferralsStatisticRequest, opts ...grpc.CallOption) (*referral.GetReferralsStatisticResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetReferralsStatistic")
	}

	var r0 *referral.GetReferralsStatisticResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *referral.GetReferralsStatisticRequest, ...grpc.CallOption) (*referral.GetReferralsStatisticResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *referral.GetReferralsStatisticRequest, ...grpc.CallOption) *referral.GetReferralsStatisticResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*referral.GetReferralsStatisticResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *referral.GetReferralsStatisticRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAccountingReferralServiceClient_GetReferralsStatistic_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetReferralsStatistic'
type MockAccountingReferralServiceClient_GetReferralsStatistic_Call struct {
	*mock.Call
}

// GetReferralsStatistic is a helper method to define mock.On call
//   - ctx context.Context
//   - in *referral.GetReferralsStatisticRequest
//   - opts ...grpc.CallOption
func (_e *MockAccountingReferralServiceClient_Expecter) GetReferralsStatistic(ctx interface{}, in interface{}, opts ...interface{}) *MockAccountingReferralServiceClient_GetReferralsStatistic_Call {
	return &MockAccountingReferralServiceClient_GetReferralsStatistic_Call{Call: _e.mock.On("GetReferralsStatistic",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAccountingReferralServiceClient_GetReferralsStatistic_Call) Run(run func(ctx context.Context, in *referral.GetReferralsStatisticRequest, opts ...grpc.CallOption)) *MockAccountingReferralServiceClient_GetReferralsStatistic_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*referral.GetReferralsStatisticRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockAccountingReferralServiceClient_GetReferralsStatistic_Call) Return(_a0 *referral.GetReferralsStatisticResponse, _a1 error) *MockAccountingReferralServiceClient_GetReferralsStatistic_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAccountingReferralServiceClient_GetReferralsStatistic_Call) RunAndReturn(run func(context.Context, *referral.GetReferralsStatisticRequest, ...grpc.CallOption) (*referral.GetReferralsStatisticResponse, error)) *MockAccountingReferralServiceClient_GetReferralsStatistic_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAccountingReferralServiceClient creates a new instance of MockAccountingReferralServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAccountingReferralServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAccountingReferralServiceClient {
	mock := &MockAccountingReferralServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

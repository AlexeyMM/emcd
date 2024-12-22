// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.0
// source: admin.proto

package adminpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	MerchantAdminService_CreateMerchant_FullMethodName = "/admin.MerchantAdminService/CreateMerchant"
)

// MerchantAdminServiceClient is the client API for MerchantAdminService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MerchantAdminServiceClient interface {
	CreateMerchant(ctx context.Context, in *CreateMerchantRequest, opts ...grpc.CallOption) (*CreateMerchantResponse, error)
}

type merchantAdminServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMerchantAdminServiceClient(cc grpc.ClientConnInterface) MerchantAdminServiceClient {
	return &merchantAdminServiceClient{cc}
}

func (c *merchantAdminServiceClient) CreateMerchant(ctx context.Context, in *CreateMerchantRequest, opts ...grpc.CallOption) (*CreateMerchantResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateMerchantResponse)
	err := c.cc.Invoke(ctx, MerchantAdminService_CreateMerchant_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MerchantAdminServiceServer is the server API for MerchantAdminService service.
// All implementations must embed UnimplementedMerchantAdminServiceServer
// for forward compatibility.
type MerchantAdminServiceServer interface {
	CreateMerchant(context.Context, *CreateMerchantRequest) (*CreateMerchantResponse, error)
	mustEmbedUnimplementedMerchantAdminServiceServer()
}

// UnimplementedMerchantAdminServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMerchantAdminServiceServer struct{}

func (UnimplementedMerchantAdminServiceServer) CreateMerchant(context.Context, *CreateMerchantRequest) (*CreateMerchantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMerchant not implemented")
}
func (UnimplementedMerchantAdminServiceServer) mustEmbedUnimplementedMerchantAdminServiceServer() {}
func (UnimplementedMerchantAdminServiceServer) testEmbeddedByValue()                              {}

// UnsafeMerchantAdminServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MerchantAdminServiceServer will
// result in compilation errors.
type UnsafeMerchantAdminServiceServer interface {
	mustEmbedUnimplementedMerchantAdminServiceServer()
}

func RegisterMerchantAdminServiceServer(s grpc.ServiceRegistrar, srv MerchantAdminServiceServer) {
	// If the following call pancis, it indicates UnimplementedMerchantAdminServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MerchantAdminService_ServiceDesc, srv)
}

func _MerchantAdminService_CreateMerchant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMerchantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MerchantAdminServiceServer).CreateMerchant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MerchantAdminService_CreateMerchant_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MerchantAdminServiceServer).CreateMerchant(ctx, req.(*CreateMerchantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MerchantAdminService_ServiceDesc is the grpc.ServiceDesc for MerchantAdminService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MerchantAdminService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "admin.MerchantAdminService",
	HandlerType: (*MerchantAdminServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMerchant",
			Handler:    _MerchantAdminService_CreateMerchant_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "admin.proto",
}

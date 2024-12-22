// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.28.2
// source: protocol/proto/b2bEndpointAdmin.proto

package b2bEndpointAdmin

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	EndpointAdminService_AddClient_FullMethodName         = "/b2bEndpointAdmin.endpointAdminService/AddClient"
	EndpointAdminService_GenerateKey_FullMethodName       = "/b2bEndpointAdmin.endpointAdminService/GenerateKey"
	EndpointAdminService_GetActiveKeys_FullMethodName     = "/b2bEndpointAdmin.endpointAdminService/GetActiveKeys"
	EndpointAdminService_DeactivateKey_FullMethodName     = "/b2bEndpointAdmin.endpointAdminService/DeactivateKey"
	EndpointAdminService_DeactivateAllKeys_FullMethodName = "/b2bEndpointAdmin.endpointAdminService/DeactivateAllKeys"
	EndpointAdminService_AddIPs_FullMethodName            = "/b2bEndpointAdmin.endpointAdminService/AddIPs"
	EndpointAdminService_GetIPs_FullMethodName            = "/b2bEndpointAdmin.endpointAdminService/GetIPs"
	EndpointAdminService_DeleteIP_FullMethodName          = "/b2bEndpointAdmin.endpointAdminService/DeleteIP"
	EndpointAdminService_DeleteAllIPs_FullMethodName      = "/b2bEndpointAdmin.endpointAdminService/DeleteAllIPs"
)

// EndpointAdminServiceClient is the client API for EndpointAdminService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EndpointAdminServiceClient interface {
	AddClient(ctx context.Context, in *AddClientRequest, opts ...grpc.CallOption) (*AddClientResponse, error)
	GenerateKey(ctx context.Context, in *GenerateKeyRequest, opts ...grpc.CallOption) (*GenerateKeyResponse, error)
	GetActiveKeys(ctx context.Context, in *GetActiveKeysRequest, opts ...grpc.CallOption) (*GetActiveKeysResponse, error)
	DeactivateKey(ctx context.Context, in *DeactivateKeyRequest, opts ...grpc.CallOption) (*DeactivateKeyResponse, error)
	DeactivateAllKeys(ctx context.Context, in *DeactivateAllKeysRequest, opts ...grpc.CallOption) (*DeactivateAllKeysResponse, error)
	AddIPs(ctx context.Context, in *AddIPsRequest, opts ...grpc.CallOption) (*AddIPsResponse, error)
	GetIPs(ctx context.Context, in *GetIPsRequest, opts ...grpc.CallOption) (*GetIPsResponse, error)
	DeleteIP(ctx context.Context, in *DeleteIPRequest, opts ...grpc.CallOption) (*DeleteIPResponse, error)
	DeleteAllIPs(ctx context.Context, in *DeleteAllIPsRequest, opts ...grpc.CallOption) (*DeleteAllIPsResponse, error)
}

type endpointAdminServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEndpointAdminServiceClient(cc grpc.ClientConnInterface) EndpointAdminServiceClient {
	return &endpointAdminServiceClient{cc}
}

func (c *endpointAdminServiceClient) AddClient(ctx context.Context, in *AddClientRequest, opts ...grpc.CallOption) (*AddClientResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddClientResponse)
	err := c.cc.Invoke(ctx, EndpointAdminService_AddClient_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointAdminServiceClient) GenerateKey(ctx context.Context, in *GenerateKeyRequest, opts ...grpc.CallOption) (*GenerateKeyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GenerateKeyResponse)
	err := c.cc.Invoke(ctx, EndpointAdminService_GenerateKey_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointAdminServiceClient) GetActiveKeys(ctx context.Context, in *GetActiveKeysRequest, opts ...grpc.CallOption) (*GetActiveKeysResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetActiveKeysResponse)
	err := c.cc.Invoke(ctx, EndpointAdminService_GetActiveKeys_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointAdminServiceClient) DeactivateKey(ctx context.Context, in *DeactivateKeyRequest, opts ...grpc.CallOption) (*DeactivateKeyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeactivateKeyResponse)
	err := c.cc.Invoke(ctx, EndpointAdminService_DeactivateKey_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointAdminServiceClient) DeactivateAllKeys(ctx context.Context, in *DeactivateAllKeysRequest, opts ...grpc.CallOption) (*DeactivateAllKeysResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeactivateAllKeysResponse)
	err := c.cc.Invoke(ctx, EndpointAdminService_DeactivateAllKeys_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointAdminServiceClient) AddIPs(ctx context.Context, in *AddIPsRequest, opts ...grpc.CallOption) (*AddIPsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddIPsResponse)
	err := c.cc.Invoke(ctx, EndpointAdminService_AddIPs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointAdminServiceClient) GetIPs(ctx context.Context, in *GetIPsRequest, opts ...grpc.CallOption) (*GetIPsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetIPsResponse)
	err := c.cc.Invoke(ctx, EndpointAdminService_GetIPs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointAdminServiceClient) DeleteIP(ctx context.Context, in *DeleteIPRequest, opts ...grpc.CallOption) (*DeleteIPResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteIPResponse)
	err := c.cc.Invoke(ctx, EndpointAdminService_DeleteIP_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointAdminServiceClient) DeleteAllIPs(ctx context.Context, in *DeleteAllIPsRequest, opts ...grpc.CallOption) (*DeleteAllIPsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteAllIPsResponse)
	err := c.cc.Invoke(ctx, EndpointAdminService_DeleteAllIPs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointAdminServiceServer is the server API for EndpointAdminService service.
// All implementations must embed UnimplementedEndpointAdminServiceServer
// for forward compatibility
type EndpointAdminServiceServer interface {
	AddClient(context.Context, *AddClientRequest) (*AddClientResponse, error)
	GenerateKey(context.Context, *GenerateKeyRequest) (*GenerateKeyResponse, error)
	GetActiveKeys(context.Context, *GetActiveKeysRequest) (*GetActiveKeysResponse, error)
	DeactivateKey(context.Context, *DeactivateKeyRequest) (*DeactivateKeyResponse, error)
	DeactivateAllKeys(context.Context, *DeactivateAllKeysRequest) (*DeactivateAllKeysResponse, error)
	AddIPs(context.Context, *AddIPsRequest) (*AddIPsResponse, error)
	GetIPs(context.Context, *GetIPsRequest) (*GetIPsResponse, error)
	DeleteIP(context.Context, *DeleteIPRequest) (*DeleteIPResponse, error)
	DeleteAllIPs(context.Context, *DeleteAllIPsRequest) (*DeleteAllIPsResponse, error)
	mustEmbedUnimplementedEndpointAdminServiceServer()
}

// UnimplementedEndpointAdminServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEndpointAdminServiceServer struct {
}

func (UnimplementedEndpointAdminServiceServer) AddClient(context.Context, *AddClientRequest) (*AddClientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddClient not implemented")
}
func (UnimplementedEndpointAdminServiceServer) GenerateKey(context.Context, *GenerateKeyRequest) (*GenerateKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateKey not implemented")
}
func (UnimplementedEndpointAdminServiceServer) GetActiveKeys(context.Context, *GetActiveKeysRequest) (*GetActiveKeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetActiveKeys not implemented")
}
func (UnimplementedEndpointAdminServiceServer) DeactivateKey(context.Context, *DeactivateKeyRequest) (*DeactivateKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeactivateKey not implemented")
}
func (UnimplementedEndpointAdminServiceServer) DeactivateAllKeys(context.Context, *DeactivateAllKeysRequest) (*DeactivateAllKeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeactivateAllKeys not implemented")
}
func (UnimplementedEndpointAdminServiceServer) AddIPs(context.Context, *AddIPsRequest) (*AddIPsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddIPs not implemented")
}
func (UnimplementedEndpointAdminServiceServer) GetIPs(context.Context, *GetIPsRequest) (*GetIPsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIPs not implemented")
}
func (UnimplementedEndpointAdminServiceServer) DeleteIP(context.Context, *DeleteIPRequest) (*DeleteIPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteIP not implemented")
}
func (UnimplementedEndpointAdminServiceServer) DeleteAllIPs(context.Context, *DeleteAllIPsRequest) (*DeleteAllIPsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAllIPs not implemented")
}
func (UnimplementedEndpointAdminServiceServer) mustEmbedUnimplementedEndpointAdminServiceServer() {}

// UnsafeEndpointAdminServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EndpointAdminServiceServer will
// result in compilation errors.
type UnsafeEndpointAdminServiceServer interface {
	mustEmbedUnimplementedEndpointAdminServiceServer()
}

func RegisterEndpointAdminServiceServer(s grpc.ServiceRegistrar, srv EndpointAdminServiceServer) {
	s.RegisterService(&EndpointAdminService_ServiceDesc, srv)
}

func _EndpointAdminService_AddClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointAdminServiceServer).AddClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EndpointAdminService_AddClient_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointAdminServiceServer).AddClient(ctx, req.(*AddClientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointAdminService_GenerateKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointAdminServiceServer).GenerateKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EndpointAdminService_GenerateKey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointAdminServiceServer).GenerateKey(ctx, req.(*GenerateKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointAdminService_GetActiveKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetActiveKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointAdminServiceServer).GetActiveKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EndpointAdminService_GetActiveKeys_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointAdminServiceServer).GetActiveKeys(ctx, req.(*GetActiveKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointAdminService_DeactivateKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeactivateKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointAdminServiceServer).DeactivateKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EndpointAdminService_DeactivateKey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointAdminServiceServer).DeactivateKey(ctx, req.(*DeactivateKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointAdminService_DeactivateAllKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeactivateAllKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointAdminServiceServer).DeactivateAllKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EndpointAdminService_DeactivateAllKeys_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointAdminServiceServer).DeactivateAllKeys(ctx, req.(*DeactivateAllKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointAdminService_AddIPs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddIPsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointAdminServiceServer).AddIPs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EndpointAdminService_AddIPs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointAdminServiceServer).AddIPs(ctx, req.(*AddIPsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointAdminService_GetIPs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetIPsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointAdminServiceServer).GetIPs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EndpointAdminService_GetIPs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointAdminServiceServer).GetIPs(ctx, req.(*GetIPsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointAdminService_DeleteIP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteIPRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointAdminServiceServer).DeleteIP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EndpointAdminService_DeleteIP_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointAdminServiceServer).DeleteIP(ctx, req.(*DeleteIPRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointAdminService_DeleteAllIPs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAllIPsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointAdminServiceServer).DeleteAllIPs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EndpointAdminService_DeleteAllIPs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointAdminServiceServer).DeleteAllIPs(ctx, req.(*DeleteAllIPsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EndpointAdminService_ServiceDesc is the grpc.ServiceDesc for EndpointAdminService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EndpointAdminService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "b2bEndpointAdmin.endpointAdminService",
	HandlerType: (*EndpointAdminServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddClient",
			Handler:    _EndpointAdminService_AddClient_Handler,
		},
		{
			MethodName: "GenerateKey",
			Handler:    _EndpointAdminService_GenerateKey_Handler,
		},
		{
			MethodName: "GetActiveKeys",
			Handler:    _EndpointAdminService_GetActiveKeys_Handler,
		},
		{
			MethodName: "DeactivateKey",
			Handler:    _EndpointAdminService_DeactivateKey_Handler,
		},
		{
			MethodName: "DeactivateAllKeys",
			Handler:    _EndpointAdminService_DeactivateAllKeys_Handler,
		},
		{
			MethodName: "AddIPs",
			Handler:    _EndpointAdminService_AddIPs_Handler,
		},
		{
			MethodName: "GetIPs",
			Handler:    _EndpointAdminService_GetIPs_Handler,
		},
		{
			MethodName: "DeleteIP",
			Handler:    _EndpointAdminService_DeleteIP_Handler,
		},
		{
			MethodName: "DeleteAllIPs",
			Handler:    _EndpointAdminService_DeleteAllIPs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protocol/proto/b2bEndpointAdmin.proto",
}

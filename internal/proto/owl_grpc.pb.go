// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: internal/proto/owl.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Auth_Register_FullMethodName = "/main.Auth/Register"
	Auth_Login_FullMethodName    = "/main.Auth/Login"
)

// AuthClient is the client API for Auth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
}

type authClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthClient(cc grpc.ClientConnInterface) AuthClient {
	return &authClient{cc}
}

func (c *authClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Auth_Register_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, Auth_Login_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServer is the server API for Auth service.
// All implementations must embed UnimplementedAuthServer
// for forward compatibility
type AuthServer interface {
	Register(context.Context, *RegisterRequest) (*emptypb.Empty, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	mustEmbedUnimplementedAuthServer()
}

// UnimplementedAuthServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServer struct {
}

func (UnimplementedAuthServer) Register(context.Context, *RegisterRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedAuthServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAuthServer) mustEmbedUnimplementedAuthServer() {}

// UnsafeAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServer will
// result in compilation errors.
type UnsafeAuthServer interface {
	mustEmbedUnimplementedAuthServer()
}

func RegisterAuthServer(s grpc.ServiceRegistrar, srv AuthServer) {
	s.RegisterService(&Auth_ServiceDesc, srv)
}

func _Auth_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Auth_ServiceDesc is the grpc.ServiceDesc for Auth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Auth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Auth_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Auth_Login_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/owl.proto",
}

const (
	Binary_SetBinary_FullMethodName      = "/main.Binary/SetBinary"
	Binary_GetBinary_FullMethodName      = "/main.Binary/GetBinary"
	Binary_GetAllBinaries_FullMethodName = "/main.Binary/GetAllBinaries"
	Binary_DeleteBinary_FullMethodName   = "/main.Binary/DeleteBinary"
)

// BinaryClient is the client API for Binary service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BinaryClient interface {
	SetBinary(ctx context.Context, in *SetBinaryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetBinary(ctx context.Context, in *GetBinaryRequest, opts ...grpc.CallOption) (*GetBinaryResponse, error)
	GetAllBinaries(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAllBinariesResponse, error)
	DeleteBinary(ctx context.Context, in *DeleteBinaryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type binaryClient struct {
	cc grpc.ClientConnInterface
}

func NewBinaryClient(cc grpc.ClientConnInterface) BinaryClient {
	return &binaryClient{cc}
}

func (c *binaryClient) SetBinary(ctx context.Context, in *SetBinaryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Binary_SetBinary_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *binaryClient) GetBinary(ctx context.Context, in *GetBinaryRequest, opts ...grpc.CallOption) (*GetBinaryResponse, error) {
	out := new(GetBinaryResponse)
	err := c.cc.Invoke(ctx, Binary_GetBinary_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *binaryClient) GetAllBinaries(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAllBinariesResponse, error) {
	out := new(GetAllBinariesResponse)
	err := c.cc.Invoke(ctx, Binary_GetAllBinaries_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *binaryClient) DeleteBinary(ctx context.Context, in *DeleteBinaryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Binary_DeleteBinary_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BinaryServer is the server API for Binary service.
// All implementations must embed UnimplementedBinaryServer
// for forward compatibility
type BinaryServer interface {
	SetBinary(context.Context, *SetBinaryRequest) (*emptypb.Empty, error)
	GetBinary(context.Context, *GetBinaryRequest) (*GetBinaryResponse, error)
	GetAllBinaries(context.Context, *emptypb.Empty) (*GetAllBinariesResponse, error)
	DeleteBinary(context.Context, *DeleteBinaryRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedBinaryServer()
}

// UnimplementedBinaryServer must be embedded to have forward compatible implementations.
type UnimplementedBinaryServer struct {
}

func (UnimplementedBinaryServer) SetBinary(context.Context, *SetBinaryRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetBinary not implemented")
}
func (UnimplementedBinaryServer) GetBinary(context.Context, *GetBinaryRequest) (*GetBinaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBinary not implemented")
}
func (UnimplementedBinaryServer) GetAllBinaries(context.Context, *emptypb.Empty) (*GetAllBinariesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllBinaries not implemented")
}
func (UnimplementedBinaryServer) DeleteBinary(context.Context, *DeleteBinaryRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBinary not implemented")
}
func (UnimplementedBinaryServer) mustEmbedUnimplementedBinaryServer() {}

// UnsafeBinaryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BinaryServer will
// result in compilation errors.
type UnsafeBinaryServer interface {
	mustEmbedUnimplementedBinaryServer()
}

func RegisterBinaryServer(s grpc.ServiceRegistrar, srv BinaryServer) {
	s.RegisterService(&Binary_ServiceDesc, srv)
}

func _Binary_SetBinary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetBinaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BinaryServer).SetBinary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Binary_SetBinary_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BinaryServer).SetBinary(ctx, req.(*SetBinaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Binary_GetBinary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBinaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BinaryServer).GetBinary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Binary_GetBinary_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BinaryServer).GetBinary(ctx, req.(*GetBinaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Binary_GetAllBinaries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BinaryServer).GetAllBinaries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Binary_GetAllBinaries_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BinaryServer).GetAllBinaries(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Binary_DeleteBinary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteBinaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BinaryServer).DeleteBinary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Binary_DeleteBinary_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BinaryServer).DeleteBinary(ctx, req.(*DeleteBinaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Binary_ServiceDesc is the grpc.ServiceDesc for Binary service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Binary_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.Binary",
	HandlerType: (*BinaryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetBinary",
			Handler:    _Binary_SetBinary_Handler,
		},
		{
			MethodName: "GetBinary",
			Handler:    _Binary_GetBinary_Handler,
		},
		{
			MethodName: "GetAllBinaries",
			Handler:    _Binary_GetAllBinaries_Handler,
		},
		{
			MethodName: "DeleteBinary",
			Handler:    _Binary_DeleteBinary_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/owl.proto",
}

const (
	Resource_SetResource_FullMethodName     = "/main.Resource/SetResource"
	Resource_GetResource_FullMethodName     = "/main.Resource/GetResource"
	Resource_GetAllResources_FullMethodName = "/main.Resource/GetAllResources"
	Resource_DeleteResource_FullMethodName  = "/main.Resource/DeleteResource"
)

// ResourceClient is the client API for Resource service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ResourceClient interface {
	SetResource(ctx context.Context, in *SetResourceRequest, opts ...grpc.CallOption) (*SetResourceResponse, error)
	GetResource(ctx context.Context, in *GetResourceRequest, opts ...grpc.CallOption) (*GetResourceResponse, error)
	GetAllResources(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAllResourcesResponse, error)
	DeleteResource(ctx context.Context, in *DeleteResourceRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type resourceClient struct {
	cc grpc.ClientConnInterface
}

func NewResourceClient(cc grpc.ClientConnInterface) ResourceClient {
	return &resourceClient{cc}
}

func (c *resourceClient) SetResource(ctx context.Context, in *SetResourceRequest, opts ...grpc.CallOption) (*SetResourceResponse, error) {
	out := new(SetResourceResponse)
	err := c.cc.Invoke(ctx, Resource_SetResource_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceClient) GetResource(ctx context.Context, in *GetResourceRequest, opts ...grpc.CallOption) (*GetResourceResponse, error) {
	out := new(GetResourceResponse)
	err := c.cc.Invoke(ctx, Resource_GetResource_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceClient) GetAllResources(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAllResourcesResponse, error) {
	out := new(GetAllResourcesResponse)
	err := c.cc.Invoke(ctx, Resource_GetAllResources_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceClient) DeleteResource(ctx context.Context, in *DeleteResourceRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Resource_DeleteResource_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceServer is the server API for Resource service.
// All implementations must embed UnimplementedResourceServer
// for forward compatibility
type ResourceServer interface {
	SetResource(context.Context, *SetResourceRequest) (*SetResourceResponse, error)
	GetResource(context.Context, *GetResourceRequest) (*GetResourceResponse, error)
	GetAllResources(context.Context, *emptypb.Empty) (*GetAllResourcesResponse, error)
	DeleteResource(context.Context, *DeleteResourceRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedResourceServer()
}

// UnimplementedResourceServer must be embedded to have forward compatible implementations.
type UnimplementedResourceServer struct {
}

func (UnimplementedResourceServer) SetResource(context.Context, *SetResourceRequest) (*SetResourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetResource not implemented")
}
func (UnimplementedResourceServer) GetResource(context.Context, *GetResourceRequest) (*GetResourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResource not implemented")
}
func (UnimplementedResourceServer) GetAllResources(context.Context, *emptypb.Empty) (*GetAllResourcesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllResources not implemented")
}
func (UnimplementedResourceServer) DeleteResource(context.Context, *DeleteResourceRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteResource not implemented")
}
func (UnimplementedResourceServer) mustEmbedUnimplementedResourceServer() {}

// UnsafeResourceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ResourceServer will
// result in compilation errors.
type UnsafeResourceServer interface {
	mustEmbedUnimplementedResourceServer()
}

func RegisterResourceServer(s grpc.ServiceRegistrar, srv ResourceServer) {
	s.RegisterService(&Resource_ServiceDesc, srv)
}

func _Resource_SetResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).SetResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Resource_SetResource_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).SetResource(ctx, req.(*SetResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Resource_GetResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).GetResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Resource_GetResource_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).GetResource(ctx, req.(*GetResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Resource_GetAllResources_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).GetAllResources(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Resource_GetAllResources_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).GetAllResources(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Resource_DeleteResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).DeleteResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Resource_DeleteResource_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).DeleteResource(ctx, req.(*DeleteResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Resource_ServiceDesc is the grpc.ServiceDesc for Resource service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Resource_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.Resource",
	HandlerType: (*ResourceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetResource",
			Handler:    _Resource_SetResource_Handler,
		},
		{
			MethodName: "GetResource",
			Handler:    _Resource_GetResource_Handler,
		},
		{
			MethodName: "GetAllResources",
			Handler:    _Resource_GetAllResources_Handler,
		},
		{
			MethodName: "DeleteResource",
			Handler:    _Resource_DeleteResource_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/owl.proto",
}

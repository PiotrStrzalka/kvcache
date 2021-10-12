// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package model

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// KvCacheClient is the client API for KvCache service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KvCacheClient interface {
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetReply, error)
	SetWithExpiry(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetReply, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
}

type kvCacheClient struct {
	cc grpc.ClientConnInterface
}

func NewKvCacheClient(cc grpc.ClientConnInterface) KvCacheClient {
	return &kvCacheClient{cc}
}

func (c *kvCacheClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetReply, error) {
	out := new(SetReply)
	err := c.cc.Invoke(ctx, "/KvCache/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kvCacheClient) SetWithExpiry(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetReply, error) {
	out := new(SetReply)
	err := c.cc.Invoke(ctx, "/KvCache/SetWithExpiry", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kvCacheClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/KvCache/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KvCacheServer is the server API for KvCache service.
// All implementations must embed UnimplementedKvCacheServer
// for forward compatibility
type KvCacheServer interface {
	Set(context.Context, *SetRequest) (*SetReply, error)
	SetWithExpiry(context.Context, *SetRequest) (*SetReply, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	mustEmbedUnimplementedKvCacheServer()
}

// UnimplementedKvCacheServer must be embedded to have forward compatible implementations.
type UnimplementedKvCacheServer struct {
}

func (UnimplementedKvCacheServer) Set(context.Context, *SetRequest) (*SetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedKvCacheServer) SetWithExpiry(context.Context, *SetRequest) (*SetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetWithExpiry not implemented")
}
func (UnimplementedKvCacheServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedKvCacheServer) mustEmbedUnimplementedKvCacheServer() {}

// UnsafeKvCacheServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KvCacheServer will
// result in compilation errors.
type UnsafeKvCacheServer interface {
	mustEmbedUnimplementedKvCacheServer()
}

func RegisterKvCacheServer(s grpc.ServiceRegistrar, srv KvCacheServer) {
	s.RegisterService(&KvCache_ServiceDesc, srv)
}

func _KvCache_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KvCacheServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/KvCache/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KvCacheServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KvCache_SetWithExpiry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KvCacheServer).SetWithExpiry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/KvCache/SetWithExpiry",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KvCacheServer).SetWithExpiry(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KvCache_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KvCacheServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/KvCache/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KvCacheServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KvCache_ServiceDesc is the grpc.ServiceDesc for KvCache service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KvCache_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "KvCache",
	HandlerType: (*KvCacheServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Set",
			Handler:    _KvCache_Set_Handler,
		},
		{
			MethodName: "SetWithExpiry",
			Handler:    _KvCache_SetWithExpiry_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _KvCache_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/model/kvcache.proto",
}

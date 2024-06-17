// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: newsfeed_core/v1/newsfeed.proto

package newsfeedv1

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
	NewsfeedService_GetNewsfeed_FullMethodName = "/newsfeed.v1.NewsfeedService/GetNewsfeed"
)

// NewsfeedServiceClient is the client API for NewsfeedService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NewsfeedServiceClient interface {
	GetNewsfeed(ctx context.Context, in *GetNewsfeedRequest, opts ...grpc.CallOption) (*GetNewsfeedResponse, error)
}

type newsfeedServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNewsfeedServiceClient(cc grpc.ClientConnInterface) NewsfeedServiceClient {
	return &newsfeedServiceClient{cc}
}

func (c *newsfeedServiceClient) GetNewsfeed(ctx context.Context, in *GetNewsfeedRequest, opts ...grpc.CallOption) (*GetNewsfeedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetNewsfeedResponse)
	err := c.cc.Invoke(ctx, NewsfeedService_GetNewsfeed_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NewsfeedServiceServer is the server API for NewsfeedService service.
// All implementations must embed UnimplementedNewsfeedServiceServer
// for forward compatibility
type NewsfeedServiceServer interface {
	GetNewsfeed(context.Context, *GetNewsfeedRequest) (*GetNewsfeedResponse, error)
	mustEmbedUnimplementedNewsfeedServiceServer()
}

// UnimplementedNewsfeedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNewsfeedServiceServer struct {
}

func (UnimplementedNewsfeedServiceServer) GetNewsfeed(context.Context, *GetNewsfeedRequest) (*GetNewsfeedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNewsfeed not implemented")
}
func (UnimplementedNewsfeedServiceServer) mustEmbedUnimplementedNewsfeedServiceServer() {}

// UnsafeNewsfeedServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NewsfeedServiceServer will
// result in compilation errors.
type UnsafeNewsfeedServiceServer interface {
	mustEmbedUnimplementedNewsfeedServiceServer()
}

func RegisterNewsfeedServiceServer(s grpc.ServiceRegistrar, srv NewsfeedServiceServer) {
	s.RegisterService(&NewsfeedService_ServiceDesc, srv)
}

func _NewsfeedService_GetNewsfeed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNewsfeedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewsfeedServiceServer).GetNewsfeed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NewsfeedService_GetNewsfeed_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewsfeedServiceServer).GetNewsfeed(ctx, req.(*GetNewsfeedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NewsfeedService_ServiceDesc is the grpc.ServiceDesc for NewsfeedService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NewsfeedService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "newsfeed.v1.NewsfeedService",
	HandlerType: (*NewsfeedServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNewsfeed",
			Handler:    _NewsfeedService_GetNewsfeed_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "newsfeed_core/v1/newsfeed.proto",
}
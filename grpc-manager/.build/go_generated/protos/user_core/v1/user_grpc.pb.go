// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: protos/user_core/v1/user.proto

package userv1

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
	AccountService_CreateAccount_FullMethodName   = "/user.v1.AccountService/CreateAccount"
	AccountService_GetAccount_FullMethodName      = "/user.v1.AccountService/GetAccount"
	AccountService_FindAccount_FullMethodName     = "/user.v1.AccountService/FindAccount"
	AccountService_Login_FullMethodName           = "/user.v1.AccountService/Login"
	AccountService_CreateSession_FullMethodName   = "/user.v1.AccountService/CreateSession"
	AccountService_GetFollowing_FullMethodName    = "/user.v1.AccountService/GetFollowing"
	AccountService_CreateFollowing_FullMethodName = "/user.v1.AccountService/CreateFollowing"
	AccountService_DeleteFollowing_FullMethodName = "/user.v1.AccountService/DeleteFollowing"
	AccountService_CheckFollowing_FullMethodName  = "/user.v1.AccountService/CheckFollowing"
)

// AccountServiceClient is the client API for AccountService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountServiceClient interface {
	CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error)
	GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error)
	FindAccount(ctx context.Context, in *FindAccountRequest, opts ...grpc.CallOption) (*FindAccountResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...grpc.CallOption) (*CreateSessionResponse, error)
	GetFollowing(ctx context.Context, in *GetFollowingRequest, opts ...grpc.CallOption) (*GetFollowingResponse, error)
	CreateFollowing(ctx context.Context, in *CheckFollowingRequest, opts ...grpc.CallOption) (*CheckFollowingResponse, error)
	DeleteFollowing(ctx context.Context, in *CheckFollowingRequest, opts ...grpc.CallOption) (*CheckFollowingResponse, error)
	CheckFollowing(ctx context.Context, in *CheckFollowingRequest, opts ...grpc.CallOption) (*CheckFollowingResponse, error)
}

type accountServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountServiceClient(cc grpc.ClientConnInterface) AccountServiceClient {
	return &accountServiceClient{cc}
}

func (c *accountServiceClient) CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateAccountResponse)
	err := c.cc.Invoke(ctx, AccountService_CreateAccount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAccountResponse)
	err := c.cc.Invoke(ctx, AccountService_GetAccount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) FindAccount(ctx context.Context, in *FindAccountRequest, opts ...grpc.CallOption) (*FindAccountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FindAccountResponse)
	err := c.cc.Invoke(ctx, AccountService_FindAccount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, AccountService_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...grpc.CallOption) (*CreateSessionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateSessionResponse)
	err := c.cc.Invoke(ctx, AccountService_CreateSession_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetFollowing(ctx context.Context, in *GetFollowingRequest, opts ...grpc.CallOption) (*GetFollowingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFollowingResponse)
	err := c.cc.Invoke(ctx, AccountService_GetFollowing_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) CreateFollowing(ctx context.Context, in *CheckFollowingRequest, opts ...grpc.CallOption) (*CheckFollowingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckFollowingResponse)
	err := c.cc.Invoke(ctx, AccountService_CreateFollowing_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) DeleteFollowing(ctx context.Context, in *CheckFollowingRequest, opts ...grpc.CallOption) (*CheckFollowingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckFollowingResponse)
	err := c.cc.Invoke(ctx, AccountService_DeleteFollowing_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) CheckFollowing(ctx context.Context, in *CheckFollowingRequest, opts ...grpc.CallOption) (*CheckFollowingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckFollowingResponse)
	err := c.cc.Invoke(ctx, AccountService_CheckFollowing_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountServiceServer is the server API for AccountService service.
// All implementations must embed UnimplementedAccountServiceServer
// for forward compatibility.
type AccountServiceServer interface {
	CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error)
	GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error)
	FindAccount(context.Context, *FindAccountRequest) (*FindAccountResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	CreateSession(context.Context, *CreateSessionRequest) (*CreateSessionResponse, error)
	GetFollowing(context.Context, *GetFollowingRequest) (*GetFollowingResponse, error)
	CreateFollowing(context.Context, *CheckFollowingRequest) (*CheckFollowingResponse, error)
	DeleteFollowing(context.Context, *CheckFollowingRequest) (*CheckFollowingResponse, error)
	CheckFollowing(context.Context, *CheckFollowingRequest) (*CheckFollowingResponse, error)
	mustEmbedUnimplementedAccountServiceServer()
}

// UnimplementedAccountServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAccountServiceServer struct{}

func (UnimplementedAccountServiceServer) CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccount not implemented")
}
func (UnimplementedAccountServiceServer) GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccount not implemented")
}
func (UnimplementedAccountServiceServer) FindAccount(context.Context, *FindAccountRequest) (*FindAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAccount not implemented")
}
func (UnimplementedAccountServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAccountServiceServer) CreateSession(context.Context, *CreateSessionRequest) (*CreateSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSession not implemented")
}
func (UnimplementedAccountServiceServer) GetFollowing(context.Context, *GetFollowingRequest) (*GetFollowingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFollowing not implemented")
}
func (UnimplementedAccountServiceServer) CreateFollowing(context.Context, *CheckFollowingRequest) (*CheckFollowingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFollowing not implemented")
}
func (UnimplementedAccountServiceServer) DeleteFollowing(context.Context, *CheckFollowingRequest) (*CheckFollowingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFollowing not implemented")
}
func (UnimplementedAccountServiceServer) CheckFollowing(context.Context, *CheckFollowingRequest) (*CheckFollowingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckFollowing not implemented")
}
func (UnimplementedAccountServiceServer) mustEmbedUnimplementedAccountServiceServer() {}
func (UnimplementedAccountServiceServer) testEmbeddedByValue()                        {}

// UnsafeAccountServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountServiceServer will
// result in compilation errors.
type UnsafeAccountServiceServer interface {
	mustEmbedUnimplementedAccountServiceServer()
}

func RegisterAccountServiceServer(s grpc.ServiceRegistrar, srv AccountServiceServer) {
	// If the following call pancis, it indicates UnimplementedAccountServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AccountService_ServiceDesc, srv)
}

func _AccountService_CreateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).CreateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_CreateAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).CreateAccount(ctx, req.(*CreateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_GetAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetAccount(ctx, req.(*GetAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_FindAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).FindAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_FindAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).FindAccount(ctx, req.(*FindAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_CreateSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).CreateSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_CreateSession_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).CreateSession(ctx, req.(*CreateSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetFollowing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFollowingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetFollowing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_GetFollowing_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetFollowing(ctx, req.(*GetFollowingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_CreateFollowing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckFollowingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).CreateFollowing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_CreateFollowing_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).CreateFollowing(ctx, req.(*CheckFollowingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_DeleteFollowing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckFollowingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).DeleteFollowing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_DeleteFollowing_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).DeleteFollowing(ctx, req.(*CheckFollowingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_CheckFollowing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckFollowingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).CheckFollowing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_CheckFollowing_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).CheckFollowing(ctx, req.(*CheckFollowingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AccountService_ServiceDesc is the grpc.ServiceDesc for AccountService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccountService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.v1.AccountService",
	HandlerType: (*AccountServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAccount",
			Handler:    _AccountService_CreateAccount_Handler,
		},
		{
			MethodName: "GetAccount",
			Handler:    _AccountService_GetAccount_Handler,
		},
		{
			MethodName: "FindAccount",
			Handler:    _AccountService_FindAccount_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _AccountService_Login_Handler,
		},
		{
			MethodName: "CreateSession",
			Handler:    _AccountService_CreateSession_Handler,
		},
		{
			MethodName: "GetFollowing",
			Handler:    _AccountService_GetFollowing_Handler,
		},
		{
			MethodName: "CreateFollowing",
			Handler:    _AccountService_CreateFollowing_Handler,
		},
		{
			MethodName: "DeleteFollowing",
			Handler:    _AccountService_DeleteFollowing_Handler,
		},
		{
			MethodName: "CheckFollowing",
			Handler:    _AccountService_CheckFollowing_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/user_core/v1/user.proto",
}

const (
	EmailTemplateService_GetEmailTemplateById_FullMethodName         = "/user.v1.EmailTemplateService/GetEmailTemplateById"
	EmailTemplateService_GetEmailTemplateByTemplateId_FullMethodName = "/user.v1.EmailTemplateService/GetEmailTemplateByTemplateId"
	EmailTemplateService_AddEmailTemplate_FullMethodName             = "/user.v1.EmailTemplateService/AddEmailTemplate"
	EmailTemplateService_UpdateEmailTemplate_FullMethodName          = "/user.v1.EmailTemplateService/UpdateEmailTemplate"
)

// EmailTemplateServiceClient is the client API for EmailTemplateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EmailTemplateServiceClient interface {
	GetEmailTemplateById(ctx context.Context, in *GetEmailTemplateRequest, opts ...grpc.CallOption) (*GetEmailTemplateResponse, error)
	GetEmailTemplateByTemplateId(ctx context.Context, in *GetEmailTemplateByTemplateIdRequest, opts ...grpc.CallOption) (*GetEmailTemplateResponse, error)
	AddEmailTemplate(ctx context.Context, in *AddEmailTemplateRequest, opts ...grpc.CallOption) (*AddEmailTemplateResponse, error)
	UpdateEmailTemplate(ctx context.Context, in *UpdateEmailTemplateRequest, opts ...grpc.CallOption) (*UpdateEmailTemplateResponse, error)
}

type emailTemplateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEmailTemplateServiceClient(cc grpc.ClientConnInterface) EmailTemplateServiceClient {
	return &emailTemplateServiceClient{cc}
}

func (c *emailTemplateServiceClient) GetEmailTemplateById(ctx context.Context, in *GetEmailTemplateRequest, opts ...grpc.CallOption) (*GetEmailTemplateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetEmailTemplateResponse)
	err := c.cc.Invoke(ctx, EmailTemplateService_GetEmailTemplateById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emailTemplateServiceClient) GetEmailTemplateByTemplateId(ctx context.Context, in *GetEmailTemplateByTemplateIdRequest, opts ...grpc.CallOption) (*GetEmailTemplateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetEmailTemplateResponse)
	err := c.cc.Invoke(ctx, EmailTemplateService_GetEmailTemplateByTemplateId_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emailTemplateServiceClient) AddEmailTemplate(ctx context.Context, in *AddEmailTemplateRequest, opts ...grpc.CallOption) (*AddEmailTemplateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddEmailTemplateResponse)
	err := c.cc.Invoke(ctx, EmailTemplateService_AddEmailTemplate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emailTemplateServiceClient) UpdateEmailTemplate(ctx context.Context, in *UpdateEmailTemplateRequest, opts ...grpc.CallOption) (*UpdateEmailTemplateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateEmailTemplateResponse)
	err := c.cc.Invoke(ctx, EmailTemplateService_UpdateEmailTemplate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmailTemplateServiceServer is the server API for EmailTemplateService service.
// All implementations must embed UnimplementedEmailTemplateServiceServer
// for forward compatibility.
type EmailTemplateServiceServer interface {
	GetEmailTemplateById(context.Context, *GetEmailTemplateRequest) (*GetEmailTemplateResponse, error)
	GetEmailTemplateByTemplateId(context.Context, *GetEmailTemplateByTemplateIdRequest) (*GetEmailTemplateResponse, error)
	AddEmailTemplate(context.Context, *AddEmailTemplateRequest) (*AddEmailTemplateResponse, error)
	UpdateEmailTemplate(context.Context, *UpdateEmailTemplateRequest) (*UpdateEmailTemplateResponse, error)
	mustEmbedUnimplementedEmailTemplateServiceServer()
}

// UnimplementedEmailTemplateServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEmailTemplateServiceServer struct{}

func (UnimplementedEmailTemplateServiceServer) GetEmailTemplateById(context.Context, *GetEmailTemplateRequest) (*GetEmailTemplateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEmailTemplateById not implemented")
}
func (UnimplementedEmailTemplateServiceServer) GetEmailTemplateByTemplateId(context.Context, *GetEmailTemplateByTemplateIdRequest) (*GetEmailTemplateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEmailTemplateByTemplateId not implemented")
}
func (UnimplementedEmailTemplateServiceServer) AddEmailTemplate(context.Context, *AddEmailTemplateRequest) (*AddEmailTemplateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddEmailTemplate not implemented")
}
func (UnimplementedEmailTemplateServiceServer) UpdateEmailTemplate(context.Context, *UpdateEmailTemplateRequest) (*UpdateEmailTemplateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEmailTemplate not implemented")
}
func (UnimplementedEmailTemplateServiceServer) mustEmbedUnimplementedEmailTemplateServiceServer() {}
func (UnimplementedEmailTemplateServiceServer) testEmbeddedByValue()                              {}

// UnsafeEmailTemplateServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EmailTemplateServiceServer will
// result in compilation errors.
type UnsafeEmailTemplateServiceServer interface {
	mustEmbedUnimplementedEmailTemplateServiceServer()
}

func RegisterEmailTemplateServiceServer(s grpc.ServiceRegistrar, srv EmailTemplateServiceServer) {
	// If the following call pancis, it indicates UnimplementedEmailTemplateServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&EmailTemplateService_ServiceDesc, srv)
}

func _EmailTemplateService_GetEmailTemplateById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEmailTemplateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailTemplateServiceServer).GetEmailTemplateById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EmailTemplateService_GetEmailTemplateById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailTemplateServiceServer).GetEmailTemplateById(ctx, req.(*GetEmailTemplateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EmailTemplateService_GetEmailTemplateByTemplateId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEmailTemplateByTemplateIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailTemplateServiceServer).GetEmailTemplateByTemplateId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EmailTemplateService_GetEmailTemplateByTemplateId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailTemplateServiceServer).GetEmailTemplateByTemplateId(ctx, req.(*GetEmailTemplateByTemplateIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EmailTemplateService_AddEmailTemplate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddEmailTemplateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailTemplateServiceServer).AddEmailTemplate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EmailTemplateService_AddEmailTemplate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailTemplateServiceServer).AddEmailTemplate(ctx, req.(*AddEmailTemplateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EmailTemplateService_UpdateEmailTemplate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEmailTemplateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailTemplateServiceServer).UpdateEmailTemplate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EmailTemplateService_UpdateEmailTemplate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailTemplateServiceServer).UpdateEmailTemplate(ctx, req.(*UpdateEmailTemplateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EmailTemplateService_ServiceDesc is the grpc.ServiceDesc for EmailTemplateService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EmailTemplateService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.v1.EmailTemplateService",
	HandlerType: (*EmailTemplateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetEmailTemplateById",
			Handler:    _EmailTemplateService_GetEmailTemplateById_Handler,
		},
		{
			MethodName: "GetEmailTemplateByTemplateId",
			Handler:    _EmailTemplateService_GetEmailTemplateByTemplateId_Handler,
		},
		{
			MethodName: "AddEmailTemplate",
			Handler:    _EmailTemplateService_AddEmailTemplate_Handler,
		},
		{
			MethodName: "UpdateEmailTemplate",
			Handler:    _EmailTemplateService_UpdateEmailTemplate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/user_core/v1/user.proto",
}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.0
// source: user.rpcFlag

package user

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

const (
	UserClient_Login_FullMethodName    = "/user.UserClient/Login"
	UserClient_UpPass_FullMethodName   = "/user.UserClient/UpPass"
	UserClient_UserAdd_FullMethodName  = "/user.UserClient/UserAdd"
	UserClient_UserEdit_FullMethodName = "/user.UserClient/UserEdit"
	UserClient_UserDel_FullMethodName  = "/user.UserClient/UserDel"
	UserClient_UserList_FullMethodName = "/user.UserClient/UserList"
)

// UserClientClient is the client API for UserClient service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserClientClient interface {
	// 用户配置
	Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error)
	UpPass(ctx context.Context, in *UpPassReq, opts ...grpc.CallOption) (*UpPassResp, error)
	UserAdd(ctx context.Context, in *UserAddReq, opts ...grpc.CallOption) (*UserResp, error)
	UserEdit(ctx context.Context, in *UserEditReq, opts ...grpc.CallOption) (*UserResp, error)
	UserDel(ctx context.Context, in *UserDelReq, opts ...grpc.CallOption) (*UserResp, error)
	UserList(ctx context.Context, in *UserListReq, opts ...grpc.CallOption) (*UserListResp, error)
}

type userClientClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClientClient(cc grpc.ClientConnInterface) UserClientClient {
	return &userClientClient{cc}
}

func (c *userClientClient) Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error) {
	out := new(LoginResp)
	err := c.cc.Invoke(ctx, UserClient_Login_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClientClient) UpPass(ctx context.Context, in *UpPassReq, opts ...grpc.CallOption) (*UpPassResp, error) {
	out := new(UpPassResp)
	err := c.cc.Invoke(ctx, UserClient_UpPass_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClientClient) UserAdd(ctx context.Context, in *UserAddReq, opts ...grpc.CallOption) (*UserResp, error) {
	out := new(UserResp)
	err := c.cc.Invoke(ctx, UserClient_UserAdd_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClientClient) UserEdit(ctx context.Context, in *UserEditReq, opts ...grpc.CallOption) (*UserResp, error) {
	out := new(UserResp)
	err := c.cc.Invoke(ctx, UserClient_UserEdit_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClientClient) UserDel(ctx context.Context, in *UserDelReq, opts ...grpc.CallOption) (*UserResp, error) {
	out := new(UserResp)
	err := c.cc.Invoke(ctx, UserClient_UserDel_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClientClient) UserList(ctx context.Context, in *UserListReq, opts ...grpc.CallOption) (*UserListResp, error) {
	out := new(UserListResp)
	err := c.cc.Invoke(ctx, UserClient_UserList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserClientServer is the server API for UserClient service.
// All implementations must embed UnimplementedUserClientServer
// for forward compatibility
type UserClientServer interface {
	// 用户配置
	Login(context.Context, *LoginReq) (*LoginResp, error)
	UpPass(context.Context, *UpPassReq) (*UpPassResp, error)
	UserAdd(context.Context, *UserAddReq) (*UserResp, error)
	UserEdit(context.Context, *UserEditReq) (*UserResp, error)
	UserDel(context.Context, *UserDelReq) (*UserResp, error)
	UserList(context.Context, *UserListReq) (*UserListResp, error)
	mustEmbedUnimplementedUserClientServer()
}

// UnimplementedUserClientServer must be embedded to have forward compatible implementations.
type UnimplementedUserClientServer struct {
}

func (UnimplementedUserClientServer) Login(context.Context, *LoginReq) (*LoginResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedUserClientServer) UpPass(context.Context, *UpPassReq) (*UpPassResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpPass not implemented")
}
func (UnimplementedUserClientServer) UserAdd(context.Context, *UserAddReq) (*UserResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserAdd not implemented")
}
func (UnimplementedUserClientServer) UserEdit(context.Context, *UserEditReq) (*UserResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserEdit not implemented")
}
func (UnimplementedUserClientServer) UserDel(context.Context, *UserDelReq) (*UserResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserDel not implemented")
}
func (UnimplementedUserClientServer) UserList(context.Context, *UserListReq) (*UserListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserList not implemented")
}
func (UnimplementedUserClientServer) mustEmbedUnimplementedUserClientServer() {}

// UnsafeUserClientServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserClientServer will
// result in compilation errors.
type UnsafeUserClientServer interface {
	mustEmbedUnimplementedUserClientServer()
}

func RegisterUserClientServer(s grpc.ServiceRegistrar, srv UserClientServer) {
	s.RegisterService(&UserClient_ServiceDesc, srv)
}

func _UserClient_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserClientServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserClient_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserClientServer).Login(ctx, req.(*LoginReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserClient_UpPass_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpPassReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserClientServer).UpPass(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserClient_UpPass_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserClientServer).UpPass(ctx, req.(*UpPassReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserClient_UserAdd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserAddReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserClientServer).UserAdd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserClient_UserAdd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserClientServer).UserAdd(ctx, req.(*UserAddReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserClient_UserEdit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserEditReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserClientServer).UserEdit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserClient_UserEdit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserClientServer).UserEdit(ctx, req.(*UserEditReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserClient_UserDel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserDelReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserClientServer).UserDel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserClient_UserDel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserClientServer).UserDel(ctx, req.(*UserDelReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserClient_UserList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserClientServer).UserList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserClient_UserList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserClientServer).UserList(ctx, req.(*UserListReq))
	}
	return interceptor(ctx, in, info, handler)
}

// UserClient_ServiceDesc is the grpc.ServiceDesc for UserClient service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserClient_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.UserClient",
	HandlerType: (*UserClientServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _UserClient_Login_Handler,
		},
		{
			MethodName: "UpPass",
			Handler:    _UserClient_UpPass_Handler,
		},
		{
			MethodName: "UserAdd",
			Handler:    _UserClient_UserAdd_Handler,
		},
		{
			MethodName: "UserEdit",
			Handler:    _UserClient_UserEdit_Handler,
		},
		{
			MethodName: "UserDel",
			Handler:    _UserClient_UserDel_Handler,
		},
		{
			MethodName: "UserList",
			Handler:    _UserClient_UserList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.rpcFlag",
}

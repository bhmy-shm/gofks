package client

import (
	"context"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/protoc/user"
	"github.com/bhmy-shm/gofks/zrpc"
	"google.golang.org/grpc"
)

type (
	LoginReq     = user.LoginReq
	LoginResp    = user.LoginResp
	PageParam    = user.PageParam
	SorterParam  = user.SorterParam
	Status       = user.Status
	UpPassReq    = user.UpPassReq
	UpPassResp   = user.UpPassResp
	User         = user.User
	UserAddReq   = user.UserAddReq
	UserDelReq   = user.UserDelReq
	UserEditReq  = user.UserEditReq
	UserListReq  = user.UserListReq
	UserListResp = user.UserListResp
	UserResp     = user.UserResp

	AccountClient interface {
		// 用户配置
		Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error)
		UpPass(ctx context.Context, in *UpPassReq, opts ...grpc.CallOption) (*UpPassResp, error)
		UserAdd(ctx context.Context, in *UserAddReq, opts ...grpc.CallOption) (*UserResp, error)
		UserEdit(ctx context.Context, in *UserEditReq, opts ...grpc.CallOption) (*UserResp, error)
		UserDel(ctx context.Context, in *UserDelReq, opts ...grpc.CallOption) (*UserResp, error)
		UserList(ctx context.Context, in *UserListReq, opts ...grpc.CallOption) (*UserListResp, error)
	}

	defaultUserClient struct {
		cli zrpc.ClientInter
	}
)

func NewUserClient(cli zrpc.ClientInter) AccountClient {
	return &defaultUserClient{
		cli: cli,
	}
}

func (m *defaultUserClient) Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error) {
	client := user.NewUserClientClient(m.cli.Conn())
	return client.Login(ctx, in, opts...)
}

func (m *defaultUserClient) UpPass(ctx context.Context, in *UpPassReq, opts ...grpc.CallOption) (*UpPassResp, error) {
	client := user.NewUserClientClient(m.cli.Conn())
	return client.UpPass(ctx, in, opts...)
}

func (m *defaultUserClient) UserAdd(ctx context.Context, in *UserAddReq, opts ...grpc.CallOption) (*UserResp, error) {
	client := user.NewUserClientClient(m.cli.Conn())
	return client.UserAdd(ctx, in, opts...)
}

func (m *defaultUserClient) UserEdit(ctx context.Context, in *UserEditReq, opts ...grpc.CallOption) (*UserResp, error) {
	client := user.NewUserClientClient(m.cli.Conn())
	return client.UserEdit(ctx, in, opts...)
}

func (m *defaultUserClient) UserDel(ctx context.Context, in *UserDelReq, opts ...grpc.CallOption) (*UserResp, error) {
	client := user.NewUserClientClient(m.cli.Conn())
	return client.UserDel(ctx, in, opts...)
}

func (m *defaultUserClient) UserList(ctx context.Context, in *UserListReq, opts ...grpc.CallOption) (*UserListResp, error) {
	client := user.NewUserClientClient(m.cli.Conn())
	return client.UserList(ctx, in, opts...)
}

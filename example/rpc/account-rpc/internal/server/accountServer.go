package server

import (
	"context"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/internal/logic"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/internal/svc"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/protoc/user"
	"github.com/bhmy-shm/gofks/zrpc"
)

type Server struct {
	svcCtx                             *svc.ServiceContext `inject:"-"`
	user.UnimplementedUserClientServer `inject:"-"`
}

func NewServer(svcCtx *svc.ServiceContext) *Server {
	return &Server{
		svcCtx: svcCtx,
	}
}

func (s *Server) RegisterServer(server *zrpc.Server) {
	user.RegisterUserClientServer(server, s)
}

func (s *Server) Login(ctx context.Context, in *user.LoginReq) (*user.LoginResp, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(in)
}

func (s *Server) UpPass(ctx context.Context, in *user.UpPassReq) (*user.UpPassResp, error) {
	l := logic.NewUpPassLogic(ctx, s.svcCtx)
	return l.UpPass(in)
}

func (s *Server) UserAdd(ctx context.Context, in *user.UserAddReq) (*user.UserResp, error) {
	l := logic.NewUserAddLogic(ctx, s.svcCtx)
	return l.UserAdd(in)
}

func (s *Server) UserEdit(ctx context.Context, in *user.UserEditReq) (*user.UserResp, error) {
	l := logic.NewUserEditLogic(ctx, s.svcCtx)
	return l.UserEdit(in)
}

func (s *Server) UserDel(ctx context.Context, in *user.UserDelReq) (*user.UserResp, error) {
	l := logic.NewUserDelLogic(ctx, s.svcCtx)
	return l.UserDel(in)
}

func (s *Server) UserList(ctx context.Context, in *user.UserListReq) (*user.UserListResp, error) {
	l := logic.NewUserListLogic(ctx, s.svcCtx)
	return l.UserList(in)
}

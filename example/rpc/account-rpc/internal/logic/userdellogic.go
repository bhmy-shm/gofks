package logic

import (
	"context"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/internal/svc"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/protoc/user"
)

type UserDelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.LoggerInter
}

func NewUserDelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDelLogic {
	return &UserDelLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		LoggerInter: logx.WithContext(ctx),
	}
}

func (l *UserDelLogic) UserDel(in *user.UserDelReq) (*user.UserResp, error) {
	// todo: add your logic here and delete this line
	return &user.UserResp{
		Result: nil,
	}, nil
}

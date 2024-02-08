package logic

import (
	"context"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/internal/svc"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/protoc/user"
)

type UserEditLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	//logx.Logger
}

func NewUserEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserEditLogic {
	return &UserEditLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		//Logger: logx.WithContext(ctx),
	}
}

func (l *UserEditLogic) UserEdit(in *user.UserEditReq) (*user.UserResp, error) {
	// todo: add your logic here and delete this line
	resp := new(user.UserResp)

	resp.Result = nil
	return resp, nil
}

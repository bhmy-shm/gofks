package logic

import (
	"context"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/internal/svc"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/protoc/user"
)

type UpPassLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	//logx.Logger
}

func NewUpPassLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpPassLogic {
	return &UpPassLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		//Logger: logx.WithContext(ctx),
	}
}

func (l *UpPassLogic) UpPass(in *user.UpPassReq) (*user.UpPassResp, error) {
	res := new(user.UpPassResp)
	return res, nil
}

package logic

import (
	"context"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/internal/svc"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/internal/types"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/protoc/user"
	"log"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.LoggerInter
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		LoggerInter: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// todo: add your logic here and delete this line
	res := new(user.LoginResp)
	res.User = new(user.User)
	res.Result = new(user.Status)

	log.Println("login start")

	res.User.Name = "shm-Text"
	res.User.Id = 12
	res.User.OrgID = 231
	res.User.RoleID = 32
	res.User.UserType = 2
	res.User.OrgName = "res.OrgName.Text"
	res.User.RoleName = "res.RoleName.Text"

	res.Result = types.SuccessStatus()

	return res, nil
}

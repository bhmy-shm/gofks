package logic

import (
	"context"
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/bhmy-shm/gofks/example/model/account"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/internal/svc"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/internal/types"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/protoc/user"
	"gorm.io/gorm"
	"strconv"
)

type UserAddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.LoggerInter
}

func NewUserAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAddLogic {
	return &UserAddLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		LoggerInter: logx.WithContext(ctx),
	}
}

func (l *UserAddLogic) UserAdd(in *user.UserAddReq) (*user.UserResp, error) {

	var (
		err          error
		isExistCount = new(int64)
		rpcResp      = new(user.UserResp)
	)

	//判断手机号是否重复
	err = l.svcCtx.UserModel().QueryCount(l.ctx, isExistCount, account.WithUserMobile(in.UserData.Mobile))
	if err != nil || *isExistCount > 0 {
		rpcResp.Result = types.InterStatus(errorx.ErrCodeDBQueryCountRepeat)
		types.SetMetaData(rpcResp.Result,
			types.AppendMD("userId", strconv.Itoa(int(in.UserData.Id))),
			types.AppendMD("mobile", in.UserData.Mobile),
		)
		return rpcResp, nil
	}

	//密码加密

	//构建新的用户
	newUser := &account.User{
		Account: in.UserData.Account,
		Name:    in.UserData.Name,
		Pass:    in.UserData.Pass,
		Gender:  in.UserData.Gender,
		Mobile:  in.UserData.Mobile,
		Phone:   in.UserData.Phone,
	}

	err = l.svcCtx.UserModel().Trans(l.ctx, func(ctx context.Context, tx *gorm.DB) error {
		return l.svcCtx.UserModel().Insert(ctx, tx, newUser)
	})

	return rpcResp, nil
}

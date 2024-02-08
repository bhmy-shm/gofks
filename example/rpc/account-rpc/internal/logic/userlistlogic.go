package logic

import (
	"github.com/bhmy-shm/gofks/core/gormx"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/bhmy-shm/gofks/example/model/account"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/internal/svc"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/internal/types"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/protoc/user"
	"golang.org/x/net/context"
)

type UserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.LoggerInter
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		LoggerInter: logx.WithContext(ctx),
	}
}

func (l *UserListLogic) UserList(in *user.UserListReq) (*user.UserListResp, error) {

	var (
		err       error
		res       = new(user.UserListResp)
		cacheKey  string
		userTotal int64
		userList  []account.User
	)
	cacheKey = "test-shm-gofk-userListKey-id"

	err = l.svcCtx.UserModel().QueryTakeScan(l.ctx, cacheKey, userList,
		func(ctx context.Context, session gormx.SqlSession) error {

			raw := `SELECT * FROM shm_test_user`
			rawArgs := []string{}

			return session.QueryScanContext(ctx, &userList, raw, rawArgs...)
		})
	if err != nil {
		logx.Error("QueryCount is failed:", err)
	}

	cacheKey = "test-shm-gofk-userListKey-total-id"

	err = l.svcCtx.UserModel().QueryTakeTotal(l.ctx, cacheKey, &userTotal,
		func(ctx context.Context, session gormx.SqlSession) error {

			raw := `SELECT COUNT(*) FROM shm_test_user`
			rawArgs := []string{}

			return session.QueryCountContext(ctx, &userTotal, raw, rawArgs...)
		})

	res.Total = userTotal
	res.UserList = parse(userList)
	res.Result = types.SuccessStatus()
	return res, err
}

func parse(list []account.User) []*user.User {

	result := make([]*user.User, len(list))
	for i, v := range list {
		result[i] = &user.User{
			Id:      uint64(v.ID),
			Account: v.Account,
			Name:    v.Name,
			Pass:    v.Pass,
			Gender:  v.Gender,
			Mobile:  v.Mobile,
			Phone:   v.Phone,
		}
	}
	return result
}

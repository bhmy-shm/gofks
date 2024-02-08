package account

import (
	"context"
	"github.com/bhmy-shm/gofks"
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/bhmy-shm/gofks/example/api/controls/types"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/protoc/user"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func (this *AccountCase) UserList(ctx *gin.Context) {

	log.Println("userList api start")

	params := &types.UserListParams{}
	response := types.UserListResponse{}
	err := ctx.ShouldBindJSON(params)
	if err != nil {
		main.ExceptionResp(ctx, errorx.ErrApiCodeShouldBindJSON)
		return
	}

	rpcReq := &user.UserListReq{
		Page: &user.PageParam{
			PageNum:  params.PageNum,
			PageSize: params.PageSize,
		},
	}

	timeCtx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	rpcResp, err := this.Ctx.AccountRpc.UserList(timeCtx, rpcReq)
	if err != nil {
		logx.Error("rpc userList Resp failed:", err)
	}

	status := rpcResp.Result
	logx.Info("rpc userList Resp status:", status)

	switch status.Code {
	case http.StatusOK, errorx.ErrCodeOK.Code():
		response.Total = rpcResp.Total
		response.List = rpcResp.UserList
		main.Successful(ctx, response)
		return
	default:
		main.InternalResp(ctx, errorx.ErrStatus(status))
		return
	}
}

func (this *AccountCase) UserDetail(ctx *gin.Context) {
	println("userdetail api start")
}

func (this *AccountCase) UserAdd(ctx *gin.Context) {
	println("userAdd api start")

	params := &types.UserAddParams{}
	err := ctx.ShouldBindJSON(params)
	if err != nil {
		main.ExceptionResp(ctx, errorx.ErrApiCodeShouldBindJSON)
		return
	}

	rpcReq := &user.UserAddReq{
		UserData: &user.User{
			Id:         params.ID,
			OrgID:      params.OrgID,
			RoleID:     params.RoleID,
			UserType:   params.UserType,
			OrgName:    params.OrgName,
			RoleName:   params.RoleName,
			Account:    params.Account,
			Name:       params.Name,
			Pass:       params.Pass,
			Gender:     params.Gender,
			Mobile:     params.Mobile,
			Phone:      params.Phone,
			LoginCount: params.LoginCount,
		},
	}

	timeCtx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	rpcResp, err := this.Ctx.AccountRpc.UserAdd(timeCtx, rpcReq)
	if err != nil {
		logx.Error("rpc userAdd Resp failed:", err)
	}

	status := rpcResp.Result
	logx.Info("rpc userAdd Resp status:", status)

	if status.Code != http.StatusOK {
		main.InternalResp(ctx, errorx.ErrStatus(status))
	} else {
		main.Successful(ctx, status)
	}

}

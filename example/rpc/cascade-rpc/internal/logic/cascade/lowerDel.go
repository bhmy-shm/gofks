package cascade

import (
	"context"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/bhmy-shm/gofks/example/rpc/cascade-rpc/internal/svc"
	"github.com/bhmy-shm/gofks/example/rpc/cascade-rpc/protoc/cascade"
)

type LowerDel struct {
	ctx    context.Context
	svcCtx *svc.CascadeContext
	logx.LoggerInter
}

func NewLowerDel(ctx context.Context, svcCtx *svc.CascadeContext) *LowerDel {
	return &LowerDel{
		ctx:         ctx,
		svcCtx:      svcCtx,
		LoggerInter: logx.WithContext(ctx),
	}
}

func (l *LowerDel) Login(in *cascade.ReqLowerDel) (*cascade.RespLowerDel, error) {
	// todo: add your logic here and delete this line
	res := new(cascade.RespLowerDel)

	return res, nil
}

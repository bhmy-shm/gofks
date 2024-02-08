package cascade

import (
	"context"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/bhmy-shm/gofks/example/rpc/cascade-rpc/internal/svc"
	"github.com/bhmy-shm/gofks/example/rpc/cascade-rpc/protoc/cascade"
)

type LowerCreate struct {
	ctx    context.Context
	svcCtx *svc.CascadeContext
	logx.LoggerInter
}

func NewLowerCreate(ctx context.Context, svcCtx *svc.CascadeContext) *LowerCreate {
	return &LowerCreate{
		ctx:         ctx,
		svcCtx:      svcCtx,
		LoggerInter: logx.WithContext(ctx),
	}
}

func (l *LowerCreate) LowerCreate(in *cascade.ReqLowerCreate) (*cascade.RespLowerCreate, error) {
	// todo: add your logic here and delete this line
	res := new(cascade.RespLowerCreate)

	return res, nil
}

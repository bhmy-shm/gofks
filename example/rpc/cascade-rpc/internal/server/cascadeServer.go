package server

import (
	"context"
	logicCascade "github.com/bhmy-shm/gofks/example/rpc/cascade-rpc/internal/logic/cascade"
	"github.com/bhmy-shm/gofks/example/rpc/cascade-rpc/internal/svc"
	"github.com/bhmy-shm/gofks/example/rpc/cascade-rpc/protoc/cascade"
	"github.com/bhmy-shm/gofks/zrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CascadeServer struct {
	svcCtx *svc.CascadeContext
	cascade.UnimplementedCascadeClientServer
}

func NewCascadeServer(svcCtx *svc.CascadeContext) *CascadeServer {
	return &CascadeServer{
		svcCtx: svcCtx,
	}
}

func (s *CascadeServer) RegisterServer(server *zrpc.Server) {
	cascade.RegisterCascadeClientServer(server, s)
}

func (s *CascadeServer) LowerCreate(ctx context.Context, in *cascade.ReqLowerCreate) (*cascade.RespLowerCreate, error) {
	l := logicCascade.NewLowerCreate(ctx, s.svcCtx)
	//return nil, status.Errorf(codes.Unimplemented, "method LowerCreate not implemented")
	return l.LowerCreate(in)
}

func (s *CascadeServer) LowerSearch(ctx context.Context, in *cascade.ReqLowerSearch) (*cascade.RespLowerSearch, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LowerSearch not implemented")
}

func (s *CascadeServer) LowerUpdate(ctx context.Context, in *cascade.ReqLowerUpdate) (*cascade.RespLowerUpdate, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LowerUpdate not implemented")
}

func (s *CascadeServer) LowerDel(ctx context.Context, in *cascade.ReqLowerDel) (*cascade.RespLowerDel, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LowerDel not implemented")
}

func (s *CascadeServer) LowerGetInfo(ctx context.Context, in *cascade.ReqLowerGetInfo) (*cascade.RespLowerGetInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LowerGetInfo not implemented")
}

func (s *CascadeServer) SuperiorGet(ctx context.Context, in *cascade.ReqSuperiorGet) (*cascade.RespSuperiorGet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SuperiorGet not implemented")
}

func (s *CascadeServer) SuperiorSet(ctx context.Context, in *cascade.ReqSuperiorSet) (*cascade.RespSuperiorSet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SuperiorSet not implemented")
}

func (s *CascadeServer) SyncTime(ctx context.Context, in *cascade.ReqSyncTime) (*cascade.RespSyncTime, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncTime not implemented")
}

func (s *CascadeServer) HeartBeat(ctx context.Context, in *cascade.ReqHeartBeat) (*cascade.RespHeartBeat, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HeartBeat not implemented")
}

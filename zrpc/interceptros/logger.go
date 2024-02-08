package interceptros

import (
	"context"
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func LoggerInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	resp, err = handler(ctx, req)
	if err != nil {
		causeErr := errorx.Cause(err)              // err类型
		if e, ok := causeErr.(*errorx.Error); ok { // 自定义错误类型
			logx.WithContext(ctx).Errorf("【RPC-SRV-ERR】 %+v", err)

			//转成grpc err
			err = status.Error(codes.Code(e.GetCode()), e.GetMessage())
		} else {
			logx.WithContext(ctx).Errorf("【RPC-SRV-ERR】 %+v", err)
		}
	}

	log.Println("--------------------------------------- logger interceptor --------------------------------------- ")
	return resp, err
}

func LoggerStreamInterceptor(ctx context.Context, req interface{}, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (resp interface{}, err error) {

	return nil, nil
}

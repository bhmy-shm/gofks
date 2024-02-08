package server

import (
	"context"
	"github.com/bhmy-shm/gofks/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
)

//通过拦截器模式提供了一种机制，以确保服务在面临不可预知的panic时，不会直接崩溃并关闭连接，而是优雅地恢复并返回一个gRPC错误给调用者。

// StreamCrashInterceptor 拦截器
func StreamCrashInterceptor(svr interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo,
	handler grpc.StreamHandler) (err error) {
	defer handleCrash(func(r interface{}) {
		err = toPanicError(r)
	})

	return handler(svr, stream)
}

// UnaryCrashInterceptor 这是一个一元RPC调用的拦截器
func UnaryCrashInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer handleCrash(func(r interface{}) {
		err = toPanicError(r)
	})

	return handler(ctx, req)
}

// 恢复panic
func handleCrash(handler func(interface{})) {
	if r := recover(); r != nil {
		handler(r)
	}
}

func toPanicError(panic interface{}) error {
	rs := string(debug.Stack())
	logx.Errorf("%v,%+v", rs, panic)
	return status.Errorf(codes.Internal, "rpc panic: %v", panic)
}

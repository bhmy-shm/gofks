package client

import (
	"context"
	"github.com/bhmy-shm/gofks/core/utils/timex"
	"google.golang.org/grpc"
	"log"
	"path"
)

// DurationInterceptor is an interceptor that logs the processing time.
func DurationInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	serverName := path.Join(cc.Target(), method)
	start := timex.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		log.Printf("fail - %s - %v - %s", serverName, req, err.Error())
	} else {
		elapsed := timex.Since(start)
		if elapsed > defaultSlowThreshold {
			log.Printf("[RPC] ok - slowcall - %s - %v - %v", serverName, req, reply)
		}
	}

	return err
}

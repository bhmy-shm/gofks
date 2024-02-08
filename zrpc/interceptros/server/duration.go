package server

import (
	"context"
	"encoding/json"
	"github.com/bhmy-shm/gofks/core/barrier"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/bhmy-shm/gofks/core/utils/timex"
	"github.com/bhmy-shm/gofks/zrpc/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"time"
)

// UnaryStatInterceptor returns a func that uses given metrics to report stats.
func UnaryStatInterceptor(metrics *metrics.Metrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := timex.Now()
		defer func() {
			duration := timex.Since(startTime)
			metrics.Add(barrier.Task{
				Duration: duration,
			})
			logDuration(ctx, info.FullMethod, req, duration)
		}()

		return handler(ctx, req)
	}
}

func logDuration(ctx context.Context, method string, req interface{}, duration time.Duration) {
	var addr string
	client, ok := peer.FromContext(ctx)
	if ok {
		addr = client.Addr.String()
	}
	content, err := json.Marshal(req)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s - %s", addr, err.Error())
	} else if duration > SlowThreshold.Load() {
		logx.WithContext(ctx).WithDuration(duration).Errorf("[RPC] slowcall - %s - %s - %s",
			addr, method, string(content))
	} else {
		logx.WithContext(ctx).WithDuration(duration).Infof("%s - %s - %s", addr, method, string(content))
	}
}

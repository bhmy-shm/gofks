package client

import (
	"github.com/bhmy-shm/gofks/core/utils/syncx"
	"google.golang.org/grpc"
	"time"
)

const defaultSlowThreshold = time.Millisecond * 500

var SlowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)

// SetSlowThreshold sets the slow threshold.
func SetSlowThreshold(threshold time.Duration) {
	SlowThreshold.Set(threshold)
}

// WithStreamClientInterceptors uses given client stream interceptors.
func WithStreamClientInterceptors(interceptors ...grpc.StreamClientInterceptor) grpc.DialOption {
	return grpc.WithChainStreamInterceptor(interceptors...)
}

// WithUnaryClientInterceptors uses given client unary interceptors.
func WithUnaryClientInterceptors(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(interceptors...)
}

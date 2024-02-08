package server

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

// WithUnaryServerInterceptors uses given server unary interceptors.
func WithUnaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(interceptors...)
}

// WithStreamServerInterceptors uses given server stream interceptors.
func WithStreamServerInterceptors(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc.ChainStreamInterceptor(interceptors...)
}

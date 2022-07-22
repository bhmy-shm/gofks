package grpc

import (
	"context"
	"crypto/tls"
	"google.golang.org/grpc"
	"log"
	"time"
)

type ClientOption func(o *clientOptions)

func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// clientOptions is gRPC Client
type clientOptions struct {
	endpoint     string
	tlsConf      *tls.Config
	timeout      time.Duration
	ints         []grpc.UnaryClientInterceptor
	grpcOpts     []grpc.DialOption
	balancerName string
	logger       log.Logger
}

// Dial returns a GRPC connection.
func DialCtx(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := clientOptions{
		timeout: 2000 * time.Millisecond,
	}

	for _, o := range opts {
		o(&options)
	}

	return grpc.DialContext(ctx, options.endpoint, grpc.WithInsecure())
}

// DialInsecure returns an insecure GRPC connection.
func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

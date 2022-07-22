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
func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := clientOptions{
		timeout: 2000 * time.Millisecond,
	}

	for _, o := range opts {
		o(&options)
	}

	//grpcOpts := []grpc.DialOption{
	//	grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, options.balancerName)),
	//	//grpc.WithChainUnaryInterceptor(ints...),
	//}
	//if insecure {
	//	grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcinsecure.NewCredentials()))
	//}
	//if options.tlsConf != nil {
	//	grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(options.tlsConf)))
	//}
	//if len(options.grpcOpts) > 0 {
	//	grpcOpts = append(grpcOpts, options.grpcOpts...)
	//}
	return grpc.DialContext(ctx, options.endpoint, grpc.WithInsecure())
}

// DialInsecure returns an insecure GRPC connection.
func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

//func unaryClientInterceptor(ms []middleware.Middleware, timeout time.Duration, filters []selector.Filter) grpc.UnaryClientInterceptor {
//	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
//		//ctx = transport.NewClientContext(ctx, &Transport{
//		//	endpoint:  cc.Target(),
//		//	operation: method,
//		//	reqHeader: headerCarrier{},
//		//	filters:   filters,
//		//})
//		//if timeout > 0 {
//		//	var cancel context.CancelFunc
//		//	ctx, cancel = context.WithTimeout(ctx, timeout)
//		//	defer cancel()
//		//}
//		//h := func(ctx context.Context, req interface{}) (interface{}, error) {
//		//	if tr, ok := transport.FromClientContext(ctx); ok {
//		//		header := tr.RequestHeader()
//		//		keys := header.Keys()
//		//		keyvals := make([]string, 0, len(keys))
//		//		for _, k := range keys {
//		//			keyvals = append(keyvals, k, header.Get(k))
//		//		}
//		//		ctx = grpcmd.AppendToOutgoingContext(ctx, keyvals...)
//		//	}
//		//	return reply, invoker(ctx, method, req, reply, cc, opts...)
//		//}
//		//if len(ms) > 0 {
//		//	h = middleware.Chain(ms...)(h)
//		//}
//		//_, err := h(ctx, req)
//		return nil
//	}
//}

package zrpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/bhmy-shm/gofks/core/logx"
	clientInterceptors "github.com/bhmy-shm/gofks/zrpc/interceptros/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type ClientOptionFunc func(o *clientOptions)

// clientOptions is gRPC Client
type clientOptions struct {
	secure       bool //是否使用（credentials）来确保 gRPC 调用的安全性
	nonBlock     bool //是否使用 非阻塞模式连接
	endpoint     string
	tlsConf      *tls.Config
	timeout      time.Duration
	unaryInts    []grpc.UnaryClientInterceptor
	grpcOpts     []grpc.DialOption
	balancerName string
	logger       log.Logger
}

// ========================================================

type (
	ClientInter interface {
		Conn() *grpc.ClientConn
	}
	client struct {
		opts *clientOptions
		conn *grpc.ClientConn
	}
)

// Conn 获取连接
func (c *client) Conn() *grpc.ClientConn {
	return c.conn
}

// NewClient returns a Client.
func NewClient(target string, opts ...ClientOptionFunc) (ClientInter, error) {
	var cli client

	svcCfg := fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, "test-client")
	balancerOpt := WithDialOption(grpc.WithDefaultServiceConfig(svcCfg))
	opts = append([]ClientOptionFunc{balancerOpt}, opts...)

	if _, err := cli.dial(target, opts...); err != nil {
		return nil, err
	}

	return &cli, nil
}

func (c *client) dial(server string, opts ...ClientOptionFunc) (*grpc.ClientConn, error) {
	options := c.buildDialOptions(opts...)

	timeCtx, cancel := context.WithTimeout(context.Background(), c.opts.timeout*time.Millisecond)
	defer cancel()

	conn, err := grpc.DialContext(timeCtx, server, options...)
	if err != nil {
		err = fmt.Errorf("rpc dial: %s, error: %s, make sure rpc service is already started",
			server, err.Error())
		return nil, err
	}
	c.conn = conn

	logx.Infof("rpc client dial server:[%s], target:[%s]", server, c.conn.Target())

	return c.Conn(), nil
}

func (c *client) buildDialOptions(opts ...ClientOptionFunc) []grpc.DialOption {

	c.opts = new(clientOptions)

	for _, fn := range opts {
		fn(c.opts)
	}

	var options []grpc.DialOption
	if !c.opts.secure {
		options = append([]grpc.DialOption(nil), grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if !c.opts.nonBlock {
		options = append(options, grpc.WithBlock())
	}

	options = append(options,
		clientInterceptors.WithUnaryClientInterceptors(
			clientInterceptors.UnaryTracingInterceptor,
			clientInterceptors.DurationInterceptor,
			clientInterceptors.TimeoutInterceptor(c.opts.timeout),
		),
		clientInterceptors.WithStreamClientInterceptors(
			clientInterceptors.StreamTracingInterceptor,
		),
	)

	return append(options, c.opts.grpcOpts...)
}

func WithEndpoint(endpoint string) ClientOptionFunc {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// WithTimeout returns a func to customize a ClientOptions with given timeout.
func WithTimeout(timeout time.Duration) ClientOptionFunc {
	return func(options *clientOptions) {
		options.timeout = timeout
	}
}

func WithNonBlock() ClientOptionFunc {
	return func(options *clientOptions) {
		options.nonBlock = true
	}
}

func WithDialOption(opt grpc.DialOption) ClientOptionFunc {
	return func(options *clientOptions) {
		options.grpcOpts = append(options.grpcOpts, opt)
	}
}

// WithTransportCredentials return a func to make the gRPC calls secured with given credentials.
func WithTransportCredentials(creds credentials.TransportCredentials) ClientOptionFunc {
	return func(options *clientOptions) {
		options.secure = true
		options.grpcOpts = append(options.grpcOpts, grpc.WithTransportCredentials(creds))
	}
}

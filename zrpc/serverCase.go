package zrpc

import (
	"context"
	"crypto/tls"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	gofkConfs "github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/bhmy-shm/gofks/core/register"
	"github.com/bhmy-shm/gofks/zrpc/interceptros/server"
	"github.com/bhmy-shm/gofks/zrpc/metrics"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/url"
	"time"
)

// serverCase 内部实现真实加载grpc 服务端

type (
	SvrOptionFunc func(option *serverOption)

	serverOption struct {
		tlsConf  *tls.Config
		lis      net.Listener
		Err      error
		name     string
		network  string
		version  string
		address  string
		timeout  time.Duration
		metrics  *metrics.Metrics
		registry *register.EtcdRegistry
	}

	Server struct {
		*grpc.Server
		Opts       *serverOption
		instance   *register.ServiceInstance
		health     *health.Server
		Endpoint   *url.URL
		unaryInts  []grpc.UnaryServerInterceptor
		streamInts []grpc.StreamServerInterceptor
		grpcOpts   []grpc.ServerOption
		Registers  []RpcRegisterInter
	}
)

func defaultServer() *Server {
	return &Server{

		//默认参数
		Opts: &serverOption{
			network: "tcp",
			address: "127.0.0.1:8083",
		},

		//默认服务
		health: health.NewServer(),
		grpcOpts: []grpc.ServerOption{
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle: time.Minute * 5,
			}),
		},
	}
}

func (a *Server) buildInstance() *register.ServiceInstance {
	return &register.ServiceInstance{
		ID:       uuid.New().String(),
		Name:     a.Opts.name,
		Version:  a.Opts.version,
		Address:  a.Opts.address,
		Metadata: nil,
	}
}

func (s *Server) listenAndEndpoint() error {
	if s.Opts.lis == nil {
		lis, err := net.Listen(s.Opts.network, s.Opts.address)
		if err != nil {
			return err
		}
		s.Opts.lis = lis
	}
	addr, err := Extract(s.Opts.address, s.Opts.lis)
	if err != nil {
		_ = s.Opts.lis.Close()
		return err
	}
	s.Endpoint = NewEndpoint("grpc", addr, s.Opts.tlsConf != nil)
	println("")
	return nil
}

func (s *Server) AddGrpcOptions(options ...grpc.ServerOption) {
	if s.grpcOpts == nil {
		s.grpcOpts = options
	}
	s.grpcOpts = append(s.grpcOpts, options...)
}

func (s *Server) AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	if s.streamInts == nil {
		s.streamInts = interceptors
	} else {
		s.streamInts = append(s.streamInts, interceptors...)
	}
}

func (s *Server) AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	if s.unaryInts == nil {
		s.unaryInts = interceptors
	} else {
		s.unaryInts = append(s.unaryInts, interceptors...)
	}
}

// Start start the gRPC server.
func (s *Server) Start() error {

	s.instance = s.buildInstance()

	//grpc options
	options := append(s.grpcOpts,
		server.WithUnaryServerInterceptors(s.unaryInts...),
		server.WithStreamServerInterceptors(s.streamInts...),
	)
	s.Server = grpc.NewServer(options...)

	//加载具体服务
	for _, inter := range s.Registers {
		inter.RegisterServer(s)
	}

	//发起监听
	s.Opts.Err = s.listenAndEndpoint()

	//注册grpc
	grpc_health_v1.RegisterHealthServer(s.Server, s.health)
	s.health.Resume()

	reflection.Register(s.Server)

	//注册etcd
	if s.Opts.registry != nil {
		rctx, rcancel := context.WithTimeout(context.Background(), s.Opts.registry.GetDialTimeout())
		defer rcancel()

		err := s.Opts.registry.Register(rctx, s.instance)
		if err != nil {
			logx.Error("服务注册出现异常：", err)
		}
	}

	log.Printf("[gRPC] server listening on: %s\n", s.Opts.lis.Addr().String())

	return s.Serve(s.Opts.lis)
}

// Stop stop the gRPC server.
func (s *Server) Stop() error {
	s.health.Shutdown()
	s.GracefulStop()

	//反向注册
	if s.Opts.registry != nil {
		ctx, cancel := context.WithTimeout(context.Background(), s.Opts.registry.GetDialTimeout())
		defer cancel()
		if err := s.Opts.registry.Deregister(ctx, s.instance); err != nil {
			logx.Error("调用Deregister反向注册失败：", err)
			return err
		}
	}

	logx.Info("[gRPC] server stopping")
	return nil
}

// ServerLoad 加载依赖项
func ServerLoad(config *gofkConf.Config) (opts []SvrOptionFunc) {

	if len(config.GetServer().Listener()) > 0 {
		opts = append(opts, network("tcp"), address(config.GetServer().Listener()))
	}

	if config.GetServer().Timeout() > 0 {
		opts = append(opts, timeout(time.Duration(config.GetServer().Timeout())*time.Second))
	}

	if config.GetServer().EnableMetrics() {
		opts = append(opts, metricsSvr(config.GetServer().Listener()))
	}

	if config.GetRegister().IsEnable() {
		opts = append(opts, etcd(config.GetRegister()))
	}

	return
}

// Network with server network.
func network(network string) SvrOptionFunc {
	return func(option *serverOption) {
		option.network = network
	}
}

// Address with server address.
func address(addr string) SvrOptionFunc {
	return func(option *serverOption) {
		option.address = addr
	}
}

// Timeout with server timeout.
func timeout(timeout time.Duration) SvrOptionFunc {
	return func(option *serverOption) {
		option.timeout = timeout
	}
}

// Metrics Must after Address or input address
func metricsSvr(address string) SvrOptionFunc {
	return func(option *serverOption) {
		option.address = address
		option.metrics = metrics.NewMetrics(option.address)
	}
}

// Etcd with Server register
func etcd(config *gofkConfs.RegistryConfig) SvrOptionFunc {
	return func(option *serverOption) {
		option.registry = register.NewEtcdRegistry(config)
	}
}

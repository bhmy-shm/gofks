package zrpc

import (
	pkgConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/zrpc/interceptros/server"
	"google.golang.org/grpc"
	"log"
	"net/url"
)

type (
	RpcInter interface {
		AddGrpcOptions(options ...grpc.ServerOption)
		AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor)
		AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor)
		Start()
		Stop()
	}

	RpcRegisterInter interface {
		RegisterServer(*Server)
	}

	RpcServer struct {
		server   *Server
		register RpcRegisterInter
	}
)

func NewRpcServer(conf *pkgConf.Config) *RpcServer {
	rpcServer, err := newRpc(conf)
	if err != nil {
		log.Fatal(err)
	}

	return rpcServer
}

func newRpc(conf *pkgConf.Config) (*RpcServer, error) {

	var (
		opts   []SvrOptionFunc
		rpcSvr = new(RpcServer)
	)

	rpcSvr.server = defaultServer()

	//设置依赖配置项
	opts = append(opts, ServerLoad(conf)...)

	//设置额外中间件
	rpcSvr.AddUnaryInterceptors(
		server.UnaryCrashInterceptor,
		server.UnaryTracingInterceptor,
	)

	rpcSvr.AddStreamInterceptors(
		server.StreamCrashInterceptor,
		server.StreamTracingInterceptor,
	)

	//加载所有配置
	for _, fn := range opts {
		fn(rpcSvr.server.Opts)
	}

	return rpcSvr, nil
}

// Start start the RpcServer.
func (rs *RpcServer) Start() {
	if err := rs.server.Start(); err != nil {
		log.Println("server start failed:", err)
		panic(err)
	}
}

// Stop stops the RpcServer.
func (rs *RpcServer) Stop() {
	if err := rs.server.Stop(); err != nil {
		log.Println("server stop failed:", err)
		panic(err)
	}
}

// AddGrpcOptions adds given options.
func (rs *RpcServer) AddGrpcOptions(options ...grpc.ServerOption) {
	rs.server.AddGrpcOptions(options...)
}

// AddStreamInterceptors adds given stream interceptors.
func (rs *RpcServer) AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	rs.server.AddStreamInterceptors(interceptors...)
}

// AddUnaryInterceptors adds given unary interceptors.
func (rs *RpcServer) AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	rs.server.AddUnaryInterceptors(interceptors...)
}

// Endpoint return url
func (rs *RpcServer) Endpoint() (*url.URL, error) {
	if rs.server.Opts.Err != nil {
		return nil, rs.server.Opts.Err
	}
	return rs.server.Endpoint, nil
}

// Register 注册
func (rs *RpcServer) Register(inters ...RpcRegisterInter) {
	rs.server.Registers = append(rs.server.Registers, inters...)
}

package grpc

import (
	"context"
	"crypto/tls"
	"github.com/bhmy-shm/gofks/transport/host"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/url"
	"time"
)

type ServerOption func(o *Server)

type Server struct {
	*grpc.Server
	baseCtx context.Context
	tlsConf *tls.Config
	lis     net.Listener
	err     error
	network string
	address string
	timeout time.Duration

	health     *health.Server
	endpoint   *url.URL
	unaryInts  []grpc.UnaryServerInterceptor
	streamInts []grpc.StreamServerInterceptor
	grpcOpts   []grpc.ServerOption
}

// Network with server network.
func Network(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

// Address with server address.
func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// Timeout with server timeout.
func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		baseCtx: context.Background(),
		network: "tcp",
		address: ":0",
		health:  health.NewServer(),
	}
	for _, o := range opts {
		o(srv)
	}

	//unaryInts := []grpc.UnaryServerInterceptor{
	//	srv.unaryServerInterceptor(),
	//}
	//streamInts := []grpc.StreamServerInterceptor{
	//	srv.streamServerInterceptor(),
	//}
	//grpcOpts := []grpc.ServerOption{
	//	grpc.ChainUnaryInterceptor(unaryInts...),
	//	grpc.ChainStreamInterceptor(streamInts...),
	//}
	//if srv.tlsConf != nil {
	//	grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.tlsConf)))
	//}
	//if len(srv.grpcOpts) > 0 {
	//	grpcOpts = append(grpcOpts, srv.grpcOpts...)
	//}

	//
	srv.Server = grpc.NewServer()
	srv.err = srv.listenAndEndpoint()

	//
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	reflection.Register(srv.Server)
	return srv
}

func (s *Server) Endpoint() (*url.URL, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}

// Start start the gRPC server.
func (s *Server) Start(ctx context.Context) error {
	//if s.err != nil {
	//	return s.err
	//}
	//s.baseCtx = ctx

	//todo log
	log.Printf("[gRPC] server listening on: %s\n", s.lis.Addr().String())
	//s.health.Resume()
	return s.Serve(s.lis)
}

// Stop stop the gRPC server.
func (s *Server) Stop(ctx context.Context) error {
	s.health.Shutdown()
	s.GracefulStop()
	log.Println("[gRPC] server stopping")
	return nil
}

func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			return err
		}
		s.lis = lis
	}
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		_ = s.lis.Close()
		return err
	}
	s.endpoint = host.NewEndpoint("grpc", addr, s.tlsConf != nil)
	return nil
}

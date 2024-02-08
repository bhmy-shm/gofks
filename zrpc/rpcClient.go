package zrpc

import (
	gofkConf "github.com/bhmy-shm/gofks/core/config/confs"
	"google.golang.org/grpc"
	"log"
	"time"
)

type RpcClient struct {
	client ClientInter
}

// NewRpcClient returns a Client, exits on any error.
func NewRpcClient(c *gofkConf.RpcClientConf) ClientInter {
	cli, err := newClient(c)
	if err != nil {
		log.Fatal(err)
	}

	return cli
}

// NewClient returns a Client.
func newClient(c *gofkConf.RpcClientConf) (ClientInter, error) {
	var opts []ClientOptionFunc

	if c.HasCredential() {
		//opts = append(opts, WithDialOption(grpc.WithPerRPCCredentials(&auth.Credential{
		//	App:   c.App,
		//	Token: c.Token,
		//})))
	}

	if c.NonBlock() {
		opts = append(opts, WithNonBlock())
	}
	if c.Timeout() > 0 {
		opts = append(opts, WithTimeout(time.Duration(5)*time.Second))
	}

	// 生成target
	target, err := c.BuildTarget()
	if err != nil {
		return nil, err
	}

	cli, err := NewClient(target, opts...)
	if err != nil {
		return nil, err
	}

	return &RpcClient{client: cli}, nil
}

// Conn returns the underlying grpc.ClientConn.
func (rc *RpcClient) Conn() *grpc.ClientConn {
	return rc.client.Conn()
}

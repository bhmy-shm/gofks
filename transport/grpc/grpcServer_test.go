package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/bhmy-shm/gofks/transport/example"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
	"time"
)

type helloServer struct {
	example.UnimplementedHelloWorldServer
}

//func (s *helloServer) SyaHello(ctx context.Context, req *example.SayHelloRequest) (resp *example.SayHelloResponse, err error) {
//	if req.Name == "error" {
//		return nil, errors.New("say hello 调用失败")
//	}
//	if req.Name == "panic" {
//		panic("server panic")
//	}
//	log.Println("收到请求准备返回消息")
//	return &example.SayHelloResponse{Message: "hello world successful"}, nil
//}

func (s *helloServer) SyaHello(ctx context.Context, req *example.SayHelloRequest) (resp *example.SayHelloResponse, err error) {
	if req.Name == "error" {
		return nil, errors.New("say hello 调用失败")
	}
	if req.Name == "panic" {
		panic("server panic")
	}
	log.Println("收到请求准备返回消息")
	return &example.SayHelloResponse{Message: "hello world successful"}, nil
}

func (s *helloServer) mustEmbedUnimplementedHelloWorldServer() {}

func TestServer2(t *testing.T) {
	ctx := context.Background()
	//ctx = context.WithValue(ctx, testKey{}, "test")

	var opts = []ServerOption{}

	opts = append(opts, Network("tcp"))
	opts = append(opts, Address(":8083"))
	opts = append(opts, Timeout(time.Second*5))

	srv := NewServer(opts...)
	example.RegisterHelloWorldServer(srv, &helloServer{})

	e, err := srv.Endpoint()
	e2, err := srv.Endpoint()
	e3, err := srv.Endpoint()
	log.Println("e=", e, e2, e3, err)
	//if err != nil || e == nil || strings.HasSuffix(e.Host, ":8087") {
	//	t.Fatal(e, err)
	//}

	err = srv.Start(ctx)
	if err != nil {
		panic(err)
	}
}

func TestServer(t *testing.T) {

	mys := grpc.NewServer()
	example.RegisterHelloWorldServer(mys, &helloServer{})
	lis, err := net.Listen("tcp", ":8083")
	log.Println(fmt.Sprintf("DBCore服务开始启动，监听端口:%d", 8087))
	if err != nil {
		log.Fatal(err)
	}
	if err = mys.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

func TestClient(t *testing.T) {

	client, err := grpc.DialContext(context.Background(), "127.0.0.1:8083", grpc.WithInsecure())
	if err != nil {
		log.Println("dial context is failed", err)
	}

	c := example.NewHelloWorldClient(client)
	resp, err := c.SyaHello(context.Background(), &example.SayHelloRequest{
		Name: "kratos-shm",
	})
	if err != nil {
		log.Println("say hello request is failed", err)
	}
	log.Println("得到结果=", resp.Message)
}

func TestClient2(t *testing.T) {

	var uHost = "127.0.0.1:8083"

	//new grpc client
	conn, err := DialCtx(context.Background(), WithEndpoint(uHost))
	if err != nil {
		log.Fatal(err)
	}

	client := example.NewHelloWorldClient(conn)
	reply, err := client.SyaHello(context.Background(), &example.SayHelloRequest{Name: "kratos-shm"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reply.Message)

	//if !reflect.DeepEqual(reply.Message, "hello world successful") {
	//	log.Println(fmt.Errorf("expect %s, got %s", "Hello kratos", reply.Message))
	//}
}

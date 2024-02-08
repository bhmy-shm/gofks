package client

import (
	"context"
	"github.com/bhmy-shm/gofks/core/wscore"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/protoc/user"
	"github.com/bhmy-shm/gofks/zrpc"
	"github.com/goccy/go-json"
	"google.golang.org/grpc"
	"log"
	"reflect"
	"strings"
)

type (
	AccountRouter interface {
		Login() wscore.IRouter
		UpPass() wscore.IRouter
		UserAdd() wscore.IRouter
		UserEdit() wscore.IRouter
		UserDel() wscore.IRouter
		UserList() wscore.IRouter
		RegisterRouter
	}

	RegisterRouter interface {
		RegisterAccountMethods()
		Handler() wscore.IMsgHandler
	}

	defaultAccountRouter struct {
		handler wscore.IMsgHandler
		cli     zrpc.ClientInter
	}
)

func NewAccountRouter(cli zrpc.ClientInter) AccountRouter {
	router := &defaultAccountRouter{
		handler: wscore.NewMsgHandler(),
		cli:     cli,
	}
	router.RegisterAccountMethods()
	return router
}

func (r *defaultAccountRouter) Login() wscore.IRouter {
	return func(ctx context.Context, request []byte, opts ...grpc.CallOption) (interface{}, error) {

		var in *LoginReq
		err := json.Unmarshal(request, &in)
		log.Println("login unmarshal err:", err)

		client := user.NewUserClientClient(r.cli.Conn())
		return client.Login(ctx, in, opts...)
	}
}

func (r *defaultAccountRouter) UpPass() wscore.IRouter {
	return func(ctx context.Context, request []byte, opts ...grpc.CallOption) (interface{}, error) {

		var in *UpPassReq
		err := json.Unmarshal(request, &in)
		log.Println("upPass unmarshal err:", err)

		client := user.NewUserClientClient(r.cli.Conn())
		return client.UpPass(ctx, in, opts...)
	}
}

func (r *defaultAccountRouter) UserAdd() wscore.IRouter {
	return func(ctx context.Context, request []byte, opts ...grpc.CallOption) (interface{}, error) {
		var in *UserAddReq
		err := json.Unmarshal(request, &in)
		log.Println("userAdd unmarshal err:", err)

		client := user.NewUserClientClient(r.cli.Conn())
		return client.UserAdd(ctx, in, opts...)
	}
}

func (r *defaultAccountRouter) UserEdit() wscore.IRouter {
	return func(ctx context.Context, request []byte, opts ...grpc.CallOption) (interface{}, error) {

		var in *UserEditReq
		err := json.Unmarshal(request, &in)
		log.Println("userEdit unmarshal err:", err)

		client := user.NewUserClientClient(r.cli.Conn())
		return client.UserEdit(ctx, in, opts...)
	}
}

func (r *defaultAccountRouter) UserDel() wscore.IRouter {
	return func(ctx context.Context, request []byte, opts ...grpc.CallOption) (interface{}, error) {
		var in *UserDelReq
		err := json.Unmarshal(request, &in)
		log.Println("userDel unmarshal err:", err)

		client := user.NewUserClientClient(r.cli.Conn())
		return client.UserDel(ctx, in, opts...)
	}
}

func (r *defaultAccountRouter) UserList() wscore.IRouter {
	return func(ctx context.Context, request []byte, opts ...grpc.CallOption) (interface{}, error) {
		var in *UserListReq
		err := json.Unmarshal(request, &in)
		log.Println("userList unmarshal err:", err)

		client := user.NewUserClientClient(r.cli.Conn())
		return client.UserList(ctx, in, opts...)
	}
}

func (r *defaultAccountRouter) Handler() wscore.IMsgHandler {
	return r.handler
}

func (r *defaultAccountRouter) RegisterAccountMethods() {

	clientType := reflect.TypeOf(r)
	clientValue := reflect.ValueOf(r)

	methodPrefix := "account"

	for i := 0; i < clientType.NumMethod(); i++ {

		methodType := clientType.Method(i)

		if methodType.Name == "RegisterAccountMethods" || methodType.Name == "Handler" {
			continue
		}

		method := methodPrefix + "." + strings.ToLower(methodType.Name)

		//注册路由到map中
		mv := clientValue.MethodByName(methodType.Name).Call([]reflect.Value{})
		routerFunc := mv[0].Interface().(wscore.IRouter)

		r.handler.AddRouter(method, routerFunc)
	}
}

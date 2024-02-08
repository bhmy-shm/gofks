package wscore

import (
	"context"
	"github.com/bhmy-shm/gofks/core/event"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

type IMsgHandler interface {

	//DoMsgHandler 调度执行路由器，调度/执行对应的Rpc消息处理方法
	DoMsgHandler(request IRequest) IResponse

	//AddRouter 添加路由器，为消息添加具体的处理业务逻辑
	AddRouter(method string, router IRouter)
}

type IRouter func(ctx context.Context, request []byte, opts ...grpc.CallOption) (interface{}, error)

type IRequest interface {
	Message() event.IMessage
	SetMessage(message event.IMessage)

	Router() IMsgHandler

	Conn() *websocket.Conn
}

type IResponse []byte

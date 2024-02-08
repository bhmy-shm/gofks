package wscore

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

type MsgHandler struct {

	// method 协议对应的 router 路由，到rpc服务
	MethodMap map[string]IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		MethodMap: make(map[string]IRouter),
	}
}

// DoMsgHandler 1.调度执行路由器，调度/执行对应的Router消息处理方法
func (m *MsgHandler) DoMsgHandler(request IRequest) IResponse {

	handler := m.MethodMap[request.Message().GetMethod()]

	if handler == nil {
		return []byte("method not found")
	}

	//设置一个ctx上下文超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	//发起rpc调用
	rpcResp, err := handler(ctx, request.Message().GetBytes())
	if err != nil {
		log.Println("doMsgHandler err:", err)
	}

	bytes, _ := json.Marshal(rpcResp)
	log.Println("doMsgHandler rpc Response:", string(bytes))
	return bytes
}

// AddRouter 2.添加路由器，为消息添加具体的处理业务逻辑
func (m *MsgHandler) AddRouter(method string, router IRouter) {
	m.MethodMap[method] = router
}

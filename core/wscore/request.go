package wscore

import (
	"github.com/bhmy-shm/gofks/core/event"
	"github.com/gorilla/websocket"
)

type (
	requestOptions struct {

		// ws客户端连接
		conn *websocket.Conn

		// ws收到的请求数据
		msg event.IMessage

		// request to handler
		handler IMsgHandler
	}

	RequestOptionFunc func(options *requestOptions)

	request struct {
		opts *requestOptions
	}
)

func defaultRequest() *request {
	return &request{
		opts: &requestOptions{
			handler: NewMsgHandler(),
		},
	}
}

func NewRequest(opts ...RequestOptionFunc) IRequest {

	req := defaultRequest()

	for _, fn := range opts {
		fn(req.opts)
	}

	return req
}

func (r *request) Message() event.IMessage {
	if r.opts.msg == nil {
		return event.NewMessage(event.AbleJsonMsg())
	}
	return r.opts.msg
}

func (r *request) SetMessage(message event.IMessage) {
	r.opts.msg = message
}

func (r *request) Router() IMsgHandler {
	return r.opts.handler
}

func (r *request) Conn() *websocket.Conn {
	return r.opts.conn
}

func WithWsConnect(conn *websocket.Conn) RequestOptionFunc {
	return func(options *requestOptions) {
		options.conn = conn
	}
}

func WithMessage(msg event.IMessage) RequestOptionFunc {
	return func(options *requestOptions) {
		options.msg = msg
	}
}

func WithHandler(handler IMsgHandler) RequestOptionFunc {
	return func(options *requestOptions) {
		options.handler = handler
	}
}

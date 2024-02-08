package server

import (
	"context"
	"github.com/bhmy-shm/gofks/core/tracex"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
	gcodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryTracingInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	ctx, span := startSpan(ctx, info.FullMethod)
	defer span.End()

	tracex.MessageReceived.Event(ctx, 1, req)
	resp, err := handler(ctx, req)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			span.SetStatus(codes.Error, s.Message())
			span.SetAttributes(tracex.StatusCodeAttr(s.Code()))
			tracex.MessageSent.Event(ctx, 1, s.Proto())
		} else {
			span.SetStatus(codes.Error, err.Error())
		}
	}

	span.SetAttributes(tracex.StatusCodeAttr(gcodes.OK))
	tracex.MessageSent.Event(ctx, 1, resp)

	return resp, nil
}

func StreamTracingInterceptor(svr interface{}, ss grpc.ServerStream,
	info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	ctx, span := startSpan(ss.Context(), info.FullMethod)
	defer span.End()

	if err := handler(svr, wrapServerStream(ctx, ss)); err != nil {
		s, ok := status.FromError(err)
		if ok {
			span.SetStatus(codes.Error, s.Message())
			span.SetAttributes(tracex.StatusCodeAttr(s.Code()))
		} else {
			span.SetStatus(codes.Error, err.Error())
		}
	}

	span.SetAttributes(tracex.StatusCodeAttr(gcodes.OK))
	return nil
}

// wrapServerStream 将给定的 grpc.ServerStream 包装到给定的 context里。
func wrapServerStream(ctx context.Context, ss grpc.ServerStream) *serverStream {
	return &serverStream{
		ServerStream: ss,
		ctx:          ctx,
	}
}

// serverStream 包装 grpc.ServerStream,
// 拦截 RecvMsg 和 SendMsg 方法的调用
type serverStream struct {
	grpc.ServerStream
	ctx               context.Context
	receivedMessageID int
	sentMessageID     int
}

func (w *serverStream) Context() context.Context {
	return w.ctx
}

func (w *serverStream) RecvMsg(m interface{}) error {

	err := w.ServerStream.RecvMsg(m)
	if err == nil {
		// 成功接收消息后，递增消息 ID
		w.receivedMessageID++

		// 在 context 的 span 中记录接收到的消息事件
		tracex.MessageReceived.Event(w.Context(), w.receivedMessageID, m)
	}

	return err
}

func (w *serverStream) SendMsg(m interface{}) error {

	err := w.ServerStream.SendMsg(m)
	// 每次尝试发送消息，递增消息 ID
	w.sentMessageID++
	// 在 context 的 span 中记录发送的消息事件，不论发送成功与否
	tracex.MessageSent.Event(w.Context(), w.sentMessageID, m)

	return err
}

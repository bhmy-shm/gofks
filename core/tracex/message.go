package tracex

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

const messageEvent = "message"

var (
	MessageSent     = messageType(RPCMessageTypeSent)
	MessageReceived = messageType(RPCMessageTypeReceived)
)

type messageType attribute.KeyValue

// Event 将 Message类型的事件属性，添加到 OpenTelemetry 的追踪 span 中
func (m messageType) Event(ctx context.Context, id int, message interface{}) {

	span := trace.SpanFromContext(ctx)
	if p, ok := message.(proto.Message); ok {
		span.AddEvent(messageEvent, trace.WithAttributes(
			attribute.KeyValue(m),
			RPCMessageIDKey.Int(id),
			RpcMessageUncompressedSizeKey.Int(proto.Size(p)),
		))
	} else {
		span.AddEvent(messageEvent, trace.WithAttributes(
			attribute.KeyValue(m),
			RPCMessageIDKey.Int(id),
		))
	}
}

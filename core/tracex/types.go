package tracex

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.22.0"
	"google.golang.org/grpc/codes"
)

// 定义了一些远程过程调用相关的语义。规范不同的追踪系统能够以一致的方式记录和理解追踪数据。

// 在创建和记录span时，为了提供关于RPC调用的详细信息，这些attribute属性会被加入到 Span 中。
// 标准化的追踪信息，更好的过滤相关信息。

const (
	GRPCStatusCodeKey = attribute.Key("rpc.grpc.status_code")

	RPCNameKey                    = attribute.Key("rpc.name")
	RPCMessageTypeKey             = attribute.Key("message.type")
	RPCMessageIDKey               = attribute.Key("message.id")
	RPCMessageCompressedSizeKey   = attribute.Key("message.compressed_size")
	RpcMessageUncompressedSizeKey = attribute.Key("message.uncompressed_size")

	SqlMethodKey = attribute.Key("sql.method")

	HttpGinMethodKey = attribute.Key("gin.Status")
)

var (
	// RPCSystemGRPC 表示远程系统是 grpc，表明了 rpc系统使用 grpc的属性。
	RPCSystemGRPC = semconv.RPCSystemKey.String("grpc")

	// RPCNameMessage 用来注明RPC调用中涉及的消息或者操作的名称。
	RPCNameMessage = RPCNameKey.String("message")

	// RPCMessageTypeSent 用来表示已经发送的RPC消息类型。表示消息已经发出。
	RPCMessageTypeSent = RPCMessageTypeKey.String("SENT")

	// RPCMessageTypeReceived 用来表示已经接收的RPC消息类型，表示消息已被接收。
	RPCMessageTypeReceived = RPCMessageTypeKey.String("RECEIVED")
)

func StatusCodeAttr(c codes.Code) attribute.KeyValue {
	return GRPCStatusCodeKey.Int64(int64(c))
}

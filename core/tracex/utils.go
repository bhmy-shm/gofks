package tracex

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc/peer"
	"net"
	"strings"
)

func SpanInfo(fullMethod, peerAddress string) (string, []attribute.KeyValue) {
	attrs := []attribute.KeyValue{semconv.RPCSystemGRPC}
	name, mAttrs := ParseFullMethod(fullMethod)
	attrs = append(attrs, mAttrs...)
	attrs = append(attrs, PeerAttr(peerAddress)...)
	return name, attrs
}

// PeerFromCtx returns the peer from ctx.
func PeerFromCtx(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok || p == nil {
		return ""
	}

	return p.Addr.String()
}

// ParseFullMethod 返回method的方法名称 和 []attribute.KeyValue 类型的属性列表
func ParseFullMethod(fullMethod string) (string, []attribute.KeyValue) {

	//去掉开头/斜杠符号
	name := strings.TrimLeft(fullMethod, "/")

	//将 name 字符串按照斜杠进行分割，得到一个字符串切片 parts，最多包含两个元素。
	parts := strings.SplitN(name, "/", 2)
	if len(parts) != 2 {
		return name, []attribute.KeyValue(nil) //不符合规格
	}

	var attrs []attribute.KeyValue
	if service := parts[0]; service != "" {
		//将 semconv.RPCServiceKey 和 service 的字符串表示形式添加到属性列表 attrs 中。
		attrs = append(attrs, semconv.RPCServiceKey.String(service))
	}

	if method := parts[1]; method != "" {
		//将 semconv.RPCMethodKey 和 method 的字符串表示形式添加到属性列表 attrs 中。
		attrs = append(attrs, semconv.RPCMethodKey.String(method))
	}

	return name, attrs
}

// PeerAttr 返回address的 []attribute.KeyValue 属性列表
func PeerAttr(address string) []attribute.KeyValue {

	//将 addr 字符串解析为主机和端口，并将结果分别赋值给 host、port 和 err 变量。
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return []attribute.KeyValue(nil)
	}

	if len(host) == 0 {
		host = "127.0.0.1"
	}

	return []attribute.KeyValue{
		semconv.NetworkPeerAddressKey.String(host), //address 键值对
		semconv.NetworkPeerPortKey.String(port),    //port 键值对
	}
}

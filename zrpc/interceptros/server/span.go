package server

import (
	"context"
	"github.com/bhmy-shm/gofks/core/tracex"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// 启动一个新的span记录关于方法调用的信息
func startSpan(ctx context.Context, method string) (context.Context, trace.Span) {
	var md metadata.MD

	//从传入的 context.Context 中提取元数据。metadata.FromIncomingContext 函数返回一个 metadata.MD 和一个布尔值，指示是否成功提取元数据。
	requestMetadata, ok := metadata.FromIncomingContext(ctx)
	if ok {
		md = requestMetadata.Copy()
	} else {
		md = metadata.MD{}
	}

	//Extract从提供的 ctx 和元数据md中提取出跨度上下文的 spanCtx 和 bags包
	//TextMapPropagator()是用来提取和注入跨度信息的接口。
	bags, spanCtx := tracex.Extract(ctx, otel.GetTextMapPropagator(), &md)

	//创建了一个新的上下文，它包含了之前提取的bags。现在这个新的上下文ctx包含了所有的背包项
	ctx = baggage.ContextWithBaggage(ctx, bags)

	//开启一个新的追踪器。并初始化追踪器名称
	tr := otel.Tracer(tracex.TraceName)
	name, attr := tracex.SpanInfo(method, tracex.PeerFromCtx(ctx))

	//真正开启调度这个追踪器。
	return tr.Start(
		trace.ContextWithRemoteSpanContext(ctx, spanCtx), //设置新的上下文
		name,
		trace.WithSpanKind(trace.SpanKindServer), //跨度的类型为服务器
		trace.WithAttributes(attr...),            //添加跨度属性
	)
}

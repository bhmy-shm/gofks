package client

import (
	"context"
	"github.com/bhmy-shm/gofks/core/tracex"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

func startSpan(ctx context.Context, method, target string) (context.Context, trace.Span) {

	var md metadata.MD
	requestMetadata, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md = requestMetadata.Copy()
	} else {
		md = metadata.MD{}
	}

	//开启一个新的追踪器。并初始化追踪器名称
	tr := otel.Tracer(tracex.TraceName)
	name, attr := tracex.SpanInfo(method, target)

	//指示这个跨度代表一个客户端请求操作
	ctx, span := tr.Start(ctx, name,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attr...),
	)

	//将追踪上下文注入到元数据中，这样远程调用接收方可以提取和继续这个追踪
	//otel.GetTextMapPropagator 返回一个传播器，用于在不同进程间传播上下文
	tracex.Inject(ctx, otel.GetTextMapPropagator(), &md)
	ctx = metadata.NewOutgoingContext(ctx, md)

	return ctx, span
}

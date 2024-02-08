package tracex

import (
	"context"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"log"
)

func TraceIdFromContext(ctx context.Context) string {
	spanCtx := sdktrace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}
	return ""
}

func SpanIdFromContext(ctx context.Context) string {
	spanCtx := sdktrace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}
	return ""
}

var _ propagation.TextMapCarrier = (*metadataTrace)(nil)

//metadataTrace结构体通过包装 grpc的 metadata.MD类型，
//并实现 grpc metadata 的 Get、Set和Keys方法，允许元数据的提取和设置。

type metadataTrace struct {
	metadata *metadata.MD
}

func (s *metadataTrace) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (s *metadataTrace) Set(key, value string) {
	s.metadata.Set(key, value)
}

func (s *metadataTrace) Keys() []string {
	out := make([]string, 0, len(*s.metadata))
	for key := range *s.metadata {
		out = append(out, key)
	}
	return out
}

func (s *metadataTrace) Range() {
	if s.metadata.Len() > 0 {
		for key, value := range *s.metadata {
			log.Printf("key:[%v],value:[%v]\n", key, value)
		}
	}
}

//propagation.TextMapPropagator 是OpenTelemetry库定义的一种传播机制，
//用于将上下文相关的数据传递给其他组件，以便它们能够共享相同的上下文信息。

// Inject 函数用于将元数据注入到给定的上下文中，以便在不同的组件之间传递。
func Inject(ctx context.Context, p propagation.TextMapPropagator, metadata *metadata.MD) {
	p.Inject(ctx, &metadataTrace{
		metadata: metadata,
	})
}

// Extract 函数用于从给定的上下文中提取元数据。
func Extract(ctx context.Context, p propagation.TextMapPropagator, metadata *metadata.MD) (
	baggage.Baggage, sdktrace.SpanContext) {

	ctx = p.Extract(ctx, &metadataTrace{
		metadata: metadata,
	})

	return baggage.FromContext(ctx), sdktrace.SpanContextFromContext(ctx)
}

package tracex

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

const spanName = "sql"

func StartDbSpan(ctx context.Context, method string) (context.Context, oteltrace.Span) {
	tracer := otel.GetTracerProvider().Tracer(TraceName)
	start, span := tracer.Start(
		ctx,
		spanName,
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
	)
	span.SetAttributes(SqlMethodKey.String(method))

	return start, span
}

func EndDbSpan(span oteltrace.Span, err error) {
	defer span.End()

	if err == nil || err == gorm.ErrRecordNotFound {
		span.SetStatus(codes.Ok, "")
		return
	}

	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
}

//func (s *recordingSpan) RecordError(err error, opts ...trace.EventOption) {
//	if s == nil || err == nil || !s.IsRecording() {
//		return
//	}
//
//	opts = append(opts, trace.WithAttributes(
//		semconv.ExceptionTypeKey.String(typeStr(err)),
//		semconv.ExceptionMessageKey.String(err.Error()),
//	))
//
//	c := trace.NewEventConfig(opts...)
//	if c.StackTrace() {
//		opts = append(opts, trace.WithAttributes(
//			semconv.ExceptionStacktraceKey.String(recordStackTrace()),
//		))
//	}
//
//	s.addEvent(semconv.ExceptionEventName, opts...)
//}

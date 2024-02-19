package tracex

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.22.0"
	"io"
	"log"
	"os"
	"testing"
	"time"
)

const (
	TraceFileName = "testFileTrace"
)

func NewResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(), resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNamespaceKey.String("myweb")))

	return r
}

func NewStdoutExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		stdouttrace.WithPrettyPrint(),
	)
}

func NewProvider(w io.Writer) *trace.TracerProvider {
	exporter, err := NewStdoutExporter(w)
	if err != nil {
		log.Println(err)
	}

	res := NewResource()
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return tp
}

func TestTraceFileTest(t *testing.T) {

	file, err := os.OpenFile("trace.txt", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//
	tp := NewProvider(file)
	ctx, span := otel.Tracer(TraceFileName).Start(context.Background(), "span-file-start")

	time.Sleep(time.Second * 2)
	span.End()

	tp.ForceFlush(ctx)
}

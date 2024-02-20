package tracex

import (
	"context"
	"fmt"
	gofkConfs "github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/bhmy-shm/gofks/core/logx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"log"
	"os"
)

const (
	ExporterJaeger = "jaeger"
	ExporterFile   = "file"
	ExporterGrpc   = "otlpgrpc"
	ExporterHttp   = "otlphttp"
)

func StartAgent(c *gofkConfs.TraceConfig) error {

	opts := []trace.TracerProviderOption{
		// Set the sampling rate based on the parent span to 100%（50%）
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(c.Trace.Sampler))),
		// Record information about this application in a Resource.
		trace.WithResource(newResource(c)),
	}

	if len(c.Trace.Endpoint) > 0 {
		exp, err := createExporter(c)
		if err != nil {
			logx.Error(err)
			return err
		}
		opts = append(opts, trace.WithBatcher(exp)) //加载exporter
	}

	tp := trace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logx.Errorf("[otel] error: %v", err)
	}))
	return nil
}

func newResource(c *gofkConfs.TraceConfig) *resource.Resource {
	r, err := resource.Merge(
		resource.Default(), resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNamespaceKey.String(c.Trace.Namespace),
			semconv.ServiceNameKey.String(c.Trace.Name),
			semconv.ServiceVersionKey.String(c.Trace.Version),
		),
	)
	if err != nil {
		log.Println("jaeger Resource setting failed:", err)
		return nil
	}
	return r
}

func createExporter(c *gofkConfs.TraceConfig) (trace.SpanExporter, error) {
	switch c.Trace.Exporter {
	case ExporterJaeger:
		return jaeger.New(
			jaeger.WithCollectorEndpoint(
				jaeger.WithEndpoint(c.Trace.Endpoint), //指定jaeger 的地址
			),
		)
	case ExporterGrpc:
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(c.Trace.Exporter),
		}
		if len(c.Trace.OtlpHeaders) > 0 {
			opts = append(opts, otlptracegrpc.WithHeaders(c.Trace.OtlpHeaders))
		}
		return otlptracegrpc.New(context.Background(), opts...)
	case ExporterHttp:
		// Not support flexible configuration now.
		opts := []otlptracehttp.Option{
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(c.Trace.Endpoint),
		}
		if len(c.Trace.OtlpHeaders) > 0 {
			opts = append(opts, otlptracehttp.WithHeaders(c.Trace.OtlpHeaders))
		}
		if len(c.Trace.OtlpHttpPath) > 0 {
			opts = append(opts, otlptracehttp.WithURLPath(c.Trace.OtlpHttpPath))
		}
		return otlptracehttp.New(context.Background(), opts...)
	case ExporterFile:
		f, err := os.OpenFile(c.Trace.Exporter, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("file exporter endpoint error: %s", err.Error())
		}
		return stdouttrace.New(stdouttrace.WithWriter(f))
	default:
		return nil, fmt.Errorf("unknown exporter: %s", c.Trace.Exporter)
	}
}

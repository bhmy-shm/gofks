package tracex

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"log"
	"testing"
)

const (
	traceIDStr = "4bf92f3577b34da6a3ce929d0e0e4736" //traceId 必须满足32位
	spanIDStr  = "00f067aa0ba902b7"                 //spanID 必须满足16位
)

var (
	traceID = mustTraceIDFromHex(traceIDStr)
	spanID  = mustSpanIDFromHex(spanIDStr)
)

// 功能函数

// 如果输入字符串的长度不等于32，则会返回一个空的TraceID和一个表示无效TraceID长度的错误
func mustTraceIDFromHex(s string) trace.TraceID {
	var (
		t   trace.TraceID
		err error
	)
	t, err = trace.TraceIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return t
}

// 如果输入字符串的长度不等于16，则会返回一个空的SpanID和一个表示无效SpanID长度的错误
func mustSpanIDFromHex(s string) (t trace.SpanID) {
	var err error
	t, err = trace.SpanIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}

type test struct {
	name        string
	traceparent string
	tracestate  string
	sc          trace.SpanContext
}

func Test_tract_Inject_metadata_Context(t *testing.T) {
	stateStr := "key1=value1,key2=value2"
	_, err := trace.ParseTraceState(stateStr)
	state, err := trace.ParseTraceState(stateStr)
	require.NoError(t, err)

	var tests []test
	tests = append(tests, test{
		name:        "not sampled",
		traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
		sc: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: traceID,
			SpanID:  spanID,
			Remote:  true,
		}),
	})

	tests = append(tests, test{
		name:        "sampled",
		traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		sc: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
			Remote:     true,
		}),
	})
	tests = append(tests, test{
		name:        "unsupported trace rpcFlag bits dropped",
		traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		sc: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: 0xff,
			Remote:     true,
		}),
	})
	tests = append(tests, test{
		name:        "with tracestate",
		traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
		tracestate:  stateStr,
		sc: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceState: state,
			Remote:     true,
		}),
	})

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	propagator := otel.GetTextMapPropagator()

	//遍历多个测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//设置
			ctx := context.Background()
			ctx = trace.ContextWithRemoteSpanContext(ctx, tt.sc)

			//设置元数据
			want := metadata.MD{}
			want.Set("traceparent-1", tt.traceparent)
			if len(tt.tracestate) > 0 {
				want.Set("tracestate", tt.tracestate)
			}

			//注入到grpc的原数据中
			md := metadata.MD{}
			Inject(ctx, propagator, &md)
			log.Println("want,md", want, md)

			//赋值给自己实现的metadata，可以正常看到注入的内容
			mm := metadataTrace{
				metadata: &md,
			}
			log.Println("mm keys", mm.Keys())

			mm.Range()
		})
	}
}

func Test_tract_Extract_metadata_Context(t *testing.T) {
	stateStr := "key1=value1,key2=value2"
	state, err := trace.ParseTraceState(stateStr)
	require.NoError(t, err)

	tests := []struct {
		name        string
		traceparent string
		tracestate  string
		sc          trace.SpanContext
	}{
		{
			name:        "not sampled",
			traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
				SpanID:  spanID,
				Remote:  true,
			}),
		},
		{
			name:        "sampled",
			traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
				Remote:     true,
			}),
		},
		{
			name:        "valid tracestate",
			traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
			tracestate:  stateStr,
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceState: state,
				Remote:     true,
			}),
		},
		{
			name:        "invalid tracestate perserves traceparent",
			traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
			tracestate:  "invalid$@#=invalid",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
				SpanID:  spanID,
				Remote:  true,
			}),
		},
		{
			name:        "future version not sampled",
			traceparent: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
				SpanID:  spanID,
				Remote:  true,
			}),
		},
		{
			name:        "future version sampled",
			traceparent: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
				Remote:     true,
			}),
		},
		{
			name:        "future version sample bit set",
			traceparent: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-09",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
				Remote:     true,
			}),
		},
		{
			name:        "future version sample bit not set",
			traceparent: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-08",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
				SpanID:  spanID,
				Remote:  true,
			}),
		},
		{
			name:        "future version additional data",
			traceparent: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00-XYZxsf09",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
				SpanID:  spanID,
				Remote:  true,
			}),
		},
		{
			name:        "B3 format ending in dash",
			traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00-",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
				SpanID:  spanID,
				Remote:  true,
			}),
		},
		{
			name:        "future version B3 format ending in dash",
			traceparent: "03-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00-",
			sc: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
				SpanID:  spanID,
				Remote:  true,
			}),
		},
	}
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))
	propagator := otel.GetTextMapPropagator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			md := metadata.MD{}
			md.Set("traceparent", tt.traceparent)
			md.Set("tracestate", tt.tracestate)
			baage, spanCtx := Extract(ctx, propagator, &md)

			log.Println("traceParent:", tt.tracestate, tt.traceparent)
			log.Println("spanCtx:", spanCtx.SpanID())
			log.Println("bagge:", baage.String())
			log.Println("")
		})
	}
}

func TestExtractInvalidTraceContext(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{
			name:   "wrong version length",
			header: "0000-00000000000000000000000000000000-0000000000000000-01",
		},
		{
			name:   "wrong trace ID length",
			header: "00-ab00000000000000000000000000000000-cd00000000000000-01",
		},
		{
			name:   "wrong span ID length",
			header: "00-ab000000000000000000000000000000-cd0000000000000000-01",
		},
		{
			name:   "wrong trace rpcFlag length",
			header: "00-ab000000000000000000000000000000-cd00000000000000-0100",
		},
		{
			name:   "bogus version",
			header: "qw-00000000000000000000000000000000-0000000000000000-01",
		},
		{
			name:   "bogus trace ID",
			header: "00-qw000000000000000000000000000000-cd00000000000000-01",
		},
		{
			name:   "bogus span ID",
			header: "00-ab000000000000000000000000000000-qw00000000000000-01",
		},
		{
			name:   "bogus trace rpcFlag",
			header: "00-ab000000000000000000000000000000-cd00000000000000-qw",
		},
		{
			name:   "upper case version",
			header: "A0-00000000000000000000000000000000-0000000000000000-01",
		},
		{
			name:   "upper case trace ID",
			header: "00-AB000000000000000000000000000000-cd00000000000000-01",
		},
		{
			name:   "upper case span ID",
			header: "00-ab000000000000000000000000000000-CD00000000000000-01",
		},
		{
			name:   "upper case trace rpcFlag",
			header: "00-ab000000000000000000000000000000-cd00000000000000-A1",
		},
		{
			name:   "zero trace ID and span ID",
			header: "00-00000000000000000000000000000000-0000000000000000-01",
		},
		{
			name:   "trace-rpcFlag unused bits set",
			header: "00-ab000000000000000000000000000000-cd00000000000000-09",
		},
		{
			name:   "missing options",
			header: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7",
		},
		{
			name:   "empty options",
			header: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-",
		},
	}
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))
	propagator := otel.GetTextMapPropagator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			md := metadata.MD{}
			md.Set("traceparent", tt.header)
			_, spanCtx := Extract(ctx, propagator, &md)

			bb, _ := spanCtx.MarshalJSON()

			log.Println("spanCtx:", string(bb))
			log.Println("md:", md.Get("traceparent"))
			log.Println("")
		})
	}
}

func TestInvalidSpanContextDropped(t *testing.T) {
	invalidSC := trace.SpanContext{}
	require.False(t, invalidSC.IsValid())
	ctx := trace.ContextWithRemoteSpanContext(context.Background(), invalidSC)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))
	propagator := otel.GetTextMapPropagator()

	md := metadata.MD{}
	Inject(ctx, propagator, &md)
	mm := &metadataTrace{
		metadata: &md,
	}
	log.Println("mm keys:", mm.Keys())
	assert.Equal(t, "", mm.Get("traceparent"), "injected invalid SpanContext")
	log.Println("")
}

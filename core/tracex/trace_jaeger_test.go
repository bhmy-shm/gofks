package tracex

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	//semconv "go.opentelemetry.io/otel/semconv/v1.22.0"
	"log"
	"testing"
)

const TraceTestName = "testGin"

var ginTp = NewJaegerProvider()

func NewJaegerResource() *resource.Resource {
	r, err := resource.Merge(
		resource.Default(), resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNamespaceKey.String("myweb"),
			semconv.ServiceNameKey.String("testGin"),
			semconv.ServiceVersionKey.String("v0.1"),
		),
	)
	if err != nil {
		log.Println("jaeger Resource setting failed:", err)
		return nil
	}
	return r
}

func NewJaegerExporter() (trace.SpanExporter, error) {
	return jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint("http://49.235.156.213:14268/api/traces"), //指定jaeger 的地址
		),
	)
}

func NewJaegerProvider() *trace.TracerProvider {
	exporter, err := NewJaegerExporter()
	if err != nil {
		log.Println(err)
	}

	res := NewJaegerResource()
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	// 设置 W3C Trace Context 传播器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp
}

func OtelMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {

		sn := c.FullPath()
		if len(sn) == 0 {
			sn = "noRoute-" + c.Request.Method
		}

		propagator := otel.GetTextMapPropagator()
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		//创建新的span
		log.Printf("middle, c.RequestHeader:%+v \n", c.Request.Header)
		ctx, span := ginTp.Tracer(TraceTestName).Start(ctx, sn)
		defer span.End()

		//延续请求context
		propagator.Inject(ctx, propagation.HeaderCarrier(c.Writer.Header()))
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		//记录http路由请求状态
		status := c.Writer.Status()
		span.SetAttributes(HttpGinMethodKey.Int(status)) //todo 另一种写法 semconv.HTTPResponseStatusCode(status)
	}
}

func RequestSubRouting(ctx *gin.Context, reqUrl string) (gin.H, error) {
	ret := gin.H{}
	u, err := url.Parse(reqUrl)
	if err != nil {
		log.Println("http Url Parse failed:", err)
		return ret, err
	}

	if u.Host == "" {
		reqUrl = "http://127.0.0.1:8086" + u.Path
	}

	req, err := http.NewRequestWithContext(ctx.Request.Context(), "GET", reqUrl, nil)
	if err != nil {
		log.Println("httpRequest failed:", err)
		return ret, err
	}

	//子路由发起调用前，写入联路追踪信息
	otel.GetTextMapPropagator().Inject(ctx.Request.Context(), propagation.HeaderCarrier(req.Header))

	// 检查 Header 是否包含 Traceparent
	log.Printf("Headers after Inject: %+v", req.Header)

	//执行http do
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("httpRequest Do failed:", err)
		return ret, err
	}
	defer resp.Body.Close()

	// 拿到结果
	b, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(b, &ret)
	if err != nil {
		log.Println("http Resp Body Unmarshal failed:", err)
		return ret, err
	}

	return ret, nil
}

func TestGinSubRoutingJaegerTrace(t *testing.T) {
	r := gin.New()

	r.Use(OtelMiddle())

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		//子路由
		score, _ := RequestSubRouting(c, "/user/score/"+id)
		info, _ := RequestSubRouting(c, "/user/info/"+id)

		c.JSON(200, gin.H{
			"info":  info,
			"score": score,
		})
	})

	r.GET("/user/score/:id", func(c *gin.Context) {
		id := c.Param("id")
		log.Println("score:", c.Request.Header)
		c.JSON(200, gin.H{"userid": c.Param("id"), "name": "user-score" + id})
	})

	r.GET("/user/info/:id", func(c *gin.Context) {
		id := c.Param("id")
		log.Println("info:", c.Request.Header)
		c.JSON(200, gin.H{"userid": c.Param("id"), "name": "user-info" + id})
	})

	r.Run(":8086")
}

func TestJaegerTrace(t *testing.T) {
	tp := NewJaegerProvider()

	ctx, span := otel.Tracer(TraceTestName).Start(context.Background(), "jaeger-test")

	time.Sleep(time.Second * 3)
	span.End()

	tp.ForceFlush(ctx)
}

func TestGinJaegerTrace(t *testing.T) {
	r := gin.New()

	r.Use(OtelMiddle())

	r.GET("/users", func(context *gin.Context) {
		time.Sleep(time.Second * 3)
		context.JSON(200, "users 用户列表")
	})

	r.GET("/orgs", func(context *gin.Context) {
		time.Sleep(time.Second * 3)
		context.JSON(200, "orgs 组织列表")
	})

	r.Run(":8086")
}

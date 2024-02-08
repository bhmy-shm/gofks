package redisx

import (
	"context"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/bhmy-shm/gofks/core/tracex"
	"github.com/bhmy-shm/gofks/core/utils/timex"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"net"
	"strings"
	"time"
)

/*

实现了对Redis操作的跟踪和监控功能。

它通过在Redis操作前后创建和结束span来计算操作的执行时间，并根据执行时间是否超过阈值来记录慢查询日志。

同时，它还使用OpenTelemetry进行分布式跟踪，将Redis操作的命令和错误信息作为span的属性记录下来。

*/

// spanName is the span name of the redis calls.
const spanName = "redis"

const startTimeKey = "startTime"

var (
	redisCmdAttributeKey  = attribute.Key("redis.cmds")
	redisDialAttributeKey = attribute.Key("redis.dial")
	durationHook          = hook{
		tracer: otel.GetTracerProvider().Tracer(spanName),
	}
)

type hook struct {
	tracer oteltrace.Tracer
}

// DialHook 在客户端与 Redis 服务器建立新连接时触发。
// 可以使用它来修改连接过程或收集有关连接的统计信息。
func (h hook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {

		var (
			err  error
			conn net.Conn
		)

		//开启链路追踪
		ctx, span := h.startDialSpan(context.WithValue(ctx, startTimeKey, timex.Now()), network, addr)
		defer func() {
			h.endSpan(span, err)
		}()

		//执行任务
		conn, err = next(ctx, network, addr)
		return conn, err
	}
}

// ProcessHook 当一个 Redis 命令（例如 GET 或 SET）被处理时，这个钩子被触发。
// 对于单个命令，这个钩子在命令发送到 Redis 服务器之前和响应被接收之后执行自定义逻辑。
func (h hook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {

		var err error

		//开启span
		ctx, span := h.startCmdSpan(context.WithValue(ctx, startTimeKey, timex.Now()), cmd)
		defer func() {
			h.endSpan(span, err)
		}()

		val := ctx.Value(startTimeKey)
		if val == nil {
			return nil
		}

		start, ok := val.(time.Duration)
		if !ok {
			return nil
		}

		//执行hook任务
		err = next(ctx, cmd)

		//计算慢查询时长并记录日志
		duration := timex.Since(start)
		if duration > slowThreshold.Load() {
			logDuration(ctx, cmd, duration)
		}

		return err
	}
}

// ProcessPipelineHook 命令（即 pipeline）被处理时，这个钩子被触发。
// 它允许在 pipeline 的命令被批量发送和接收响应之前后注入自定义代码。
func (h hook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmdList []redis.Cmder) error {

		var batchError errorx.BatchError

		//开启span
		ctx, span := h.startCmdSpan(context.WithValue(ctx, startTimeKey, timex.Now()), cmdList...)
		defer func() {
			h.endSpan(span, batchError.Err())
		}()

		//记录所有的pipeline操作错误
		for _, cmd := range cmdList {
			if err := cmd.Err(); err == nil {
				continue
			} else {
				batchError.Add(err)
			}
		}
		val := ctx.Value(startTimeKey)
		if val == nil {
			return nil
		}

		start, ok := val.(time.Duration)
		if !ok {
			return nil
		}

		//记录慢查询日志
		duration := timex.Since(start)
		if duration > slowThreshold.Load()*time.Duration(len(cmdList)) {
			logDuration(ctx, cmdList[0], duration)
		}

		return batchError.Err()
	}
}

// ----------------------- 链路追踪的内部方法 ------------------------

// 开启cmd 的 hook的链路追踪功能
func (h hook) startCmdSpan(ctx context.Context, cmds ...redis.Cmder) (context.Context, oteltrace.Span) {
	ctx, span := h.tracer.Start(ctx,
		spanName,
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
	)

	cmdStr := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		cmdStr = append(cmdStr, cmd.Name())
	}
	span.SetAttributes(redisCmdAttributeKey.StringSlice(cmdStr))

	return ctx, span
}

func (h hook) startDialSpan(ctx context.Context, network, addr string) (context.Context, oteltrace.Span) {

	attr := tracex.PeerAttr(addr)
	attr = append(attr, redisDialAttributeKey.String(network))

	ctx, span := h.tracer.Start(ctx,
		spanName,
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
		oteltrace.WithAttributes(attr...),
	)

	return ctx, span
}

func (h hook) endSpan(span oteltrace.Span, err error) {
	defer span.End()

	if err == nil || err == gorm.ErrRecordNotFound {
		span.SetStatus(codes.Ok, "")
		return
	}
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
}

func logDuration(ctx context.Context, cmd redis.Cmder, duration time.Duration) {
	var buf strings.Builder
	for i, arg := range cmd.Args() {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(gofkConf.StrVal(arg))
	}
	logx.WithContext(ctx).WithDuration(duration).Slowf("[REDIS] slowCall on executing: %s", buf.String())
}

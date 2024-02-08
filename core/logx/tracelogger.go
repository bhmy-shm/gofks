package logx

import (
	"context"
	"fmt"
	"github.com/bhmy-shm/gofks/core/tracex"
	"github.com/bhmy-shm/gofks/core/utils/timex"
	"time"
)

func WithContext(ctx context.Context) LoggerInter {
	return &traceLogger{
		ctx: ctx,
	}
}

type traceLogger struct {
	logEntry
	ctx context.Context
}

func (l *traceLogger) Error(v ...interface{}) {
	l.err(fmt.Sprint(v...))
}

func (l *traceLogger) Errorf(format string, v ...interface{}) {
	l.err(fmt.Sprintf(format, v...))
}

func (l *traceLogger) Errorw(msg string, fields ...LogField) {
	l.err(msg, fields...)
}

func (l *traceLogger) Info(v ...interface{}) {
	l.info(fmt.Sprint(v...))
}

func (l *traceLogger) Infof(format string, v ...interface{}) {
	l.info(fmt.Sprintf(format, v...))
}

func (l *traceLogger) Infow(msg string, fields ...LogField) {
	l.info(msg, fields...)
}

func (l *traceLogger) Slow(msg ...interface{}) {
	l.slow(fmt.Sprint(msg...))
}

func (l *traceLogger) Slowf(format string, v ...interface{}) {
	l.slow(fmt.Sprintf(format, v...))
}

func (l *traceLogger) Sloww(msg string, fields ...LogField) {
	l.slow(msg, fields...)
}

func (l *traceLogger) WithContext(ctx context.Context) LoggerInter {
	if ctx == nil {
		return l
	}
	l.ctx = ctx
	return l
}

func (l *traceLogger) WithDuration(duration time.Duration) LoggerInter {
	l.Duration = timex.ReprOfDuration(duration)
	return l
}

//

func (l *traceLogger) buildFields(fields ...LogField) []LogField {
	if len(l.Duration) > 0 {
		fields = append(fields, Field(durationKey, l.Duration))
	}

	traceID := tracex.TraceIdFromContext(l.ctx)
	if len(traceID) > 0 {
		fields = append(fields, Field(traceKey, traceID))
	}

	spanID := tracex.SpanIdFromContext(l.ctx)
	if len(spanID) > 0 {
		fields = append(fields, Field(spanKey, spanID))
	}

	return fields
}

func (l *traceLogger) err(v interface{}, fields ...LogField) {
	if shallLog(ErrorLevel) {
		getWriter().Error(v, l.buildFields(fields...)...)
	}
}

func (l *traceLogger) info(v interface{}, fields ...LogField) {
	if shallLog(InfoLevel) {
		getWriter().Info(v, l.buildFields(fields...)...)
	}
}

func (l *traceLogger) stat(v interface{}, field ...LogField) {
	if shallLog(StatLevel) {
		getWriter().Stat(v, l.buildFields(field...)...)
	}
}

func (l *traceLogger) slow(v interface{}, field ...LogField) {
	if shallLog(ErrorLevel) {
		getWriter().Slow(v, l.buildFields(field...)...)
	}
}

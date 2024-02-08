package logx

import (
	"context"
	"time"
)

type LoggerInter interface {
	Error(...interface{})
	Errorf(string, ...interface{})
	Errorw(string, ...LogField)

	Info(...interface{})
	Infof(string, ...interface{})
	Infow(string, ...LogField)

	Slow(...interface{})
	Slowf(string, ...interface{})
	Sloww(string, ...LogField)

	WithContext(ctx context.Context) LoggerInter
	WithDuration(time.Duration) LoggerInter
}

package breaker

import (
	"errors"
	"fmt"
	"github.com/bhmy-shm/gofks/core/utils/proc"
)

type loggedThrottle struct {
	name             string
	internalThrottle              //内部熔断器接口，指向googleBreaker
	errWin           *errorWindow //错误窗口
}

func newLoggedThrottle(name string, t internalThrottle) loggedThrottle {
	return loggedThrottle{
		name:             name,
		internalThrottle: t,
		errWin:           new(errorWindow),
	}
}

// 尝试获取正确的 googlePromise 对象，如果成功则返回一个 PromiseWithReason。无论是否成功都会记录ErrWindow日志
func (lt loggedThrottle) allow() (Promise, error) {
	promise, err := lt.internalThrottle.allow()
	return promiseWithReason{
		promise: promise,
		errWin:  lt.errWin,
	}, lt.logError(err)
}

// 调用 internalThrottle 来执行执行实际的请求，并根据结果进行统计。
func (lt loggedThrottle) doReq(req func() error, fallback Fallback, acceptable Acceptable) error {

	err := lt.internalThrottle.doReq(req, fallback, func(err error) bool {

		//如果accept 返回false，表示错误不可以被接收记录为熔断错误，错误会被添加到错误窗口 errWin 中。（包含请求执行的错误结果）
		accept := acceptable(err)
		if !accept && err != nil {
			//记录错误日志
			lt.errWin.add(err.Error())
		}
		return accept
	})

	return lt.logError(err)
}

func (lt loggedThrottle) logError(err error) error {
	if errors.Is(err, ErrServiceUnavailable) {
		// if circuit open, not possible to have empty error window
		fmt.Sprintf(
			"proc(%s/%d), callee: %s, breaker is open and requests dropped\nlast errors:\n%s",
			proc.ProcessName(), proc.Pid(), lt.name, lt.errWin)
	}

	return err
}

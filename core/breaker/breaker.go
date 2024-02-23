package breaker

import "github.com/bhmy-shm/gofks/core/utils/random"

type Promise interface {
	// Accept 通知熔断器调用成功
	Accept()
	// Reject 通知熔断器调用失败
	Reject(reason string)
}

type internalPromise interface {
	Accept()
	Reject()
}

type (
	//Fallback 用于定义在请求被熔断器拒绝时执行的后备操作
	Fallback func(err error) error

	//Acceptable 用于判断即使在出现错误时，请求是否还是被视为成功的。
	Acceptable func(err error) bool
)

type (
	Breaker interface {
		// Name 返回熔断器的名称
		Name() string

		// Allow 检查是否允许请求。
		// 如果允许，将返回一个promise，调用者需要在成功时调用promise.Accept()，
		// 或者在失败时调用promise.Reject()。
		// 如果不允许，将返回ErrServiceUnavailable错误。
		Allow() (Promise, error)

		// Do 如果熔断器接受请求，则运行给定的请求。
		// 如果熔断器拒绝请求，Do 会立即返回错误。
		// 如果请求中发生了panic，熔断器会将其视为一个错误，并再次引发相同的panic。
		Do(req func() error) error

		// DoWithAcceptable 如果熔断器接受请求，则运行给定的请求。
		// DoWithAcceptable 会立即返回错误，如果熔断器拒绝了请求
		// 如果请求中发生了panic，熔断器会将其视为一个错误，并再次引发相同的panic。
		// acceptable 函数用来检查即使错误不为nil，调用是否还是成功的。
		DoWithAcceptable(req func() error, acceptable Acceptable) error

		// DoWithFallback 如果熔断器接受请求，则运行给定的请求。
		// DoWithFallback 如果熔断器拒绝请求，则运行fallback。
		// 如果请求中发生了panic，熔断器会将其视为一个错误，并再次引发相同的panic。
		DoWithFallback(req func() error, fallback Fallback) error

		// DoWithFallbackAcceptable 如果熔断器接受请求，则运行给定的请求。
		// DoWithFallbackAcceptable 如果熔断器拒绝请求，则运行fallback。
		// 如果请求中发生了panic，熔断器会将其视为一个错误，并再次引发相同的panic。
		// acceptable 函数用来检查即使错误不为nil，调用是否还是成功的。
		DoWithFallbackAcceptable(req func() error, fallback Fallback, acceptable Acceptable) error
	}

	Option func(breaker *circuitBreaker)

	circuitBreaker struct {
		name string
		throttle
	}

	throttle interface {
		allow() (Promise, error)
		doReq(req func() error, fallback Fallback, acceptable Acceptable) error
	}

	internalThrottle interface {
		allow() (internalPromise, error)
		doReq(req func() error, fallback Fallback, acceptable Acceptable) error
	}
)

// NewBreaker returns a Breaker object.
// opts can be used to customize the Breaker.
func NewBreaker(opts ...Option) Breaker {
	var b circuitBreaker
	for _, opt := range opts {
		opt(&b)
	}
	if len(b.name) == 0 {

		b.name = random.Rand()
	}
	b.throttle = newLoggedThrottle(b.name, newGoogleBreaker())

	return &b
}

func defaultAcceptable(err error) bool {
	return err == nil
}

func (cb *circuitBreaker) Allow() (Promise, error) {
	return cb.throttle.allow()
}

func (cb *circuitBreaker) Do(req func() error) error {
	return cb.throttle.doReq(req, nil, defaultAcceptable)
}

func (cb *circuitBreaker) DoWithAcceptable(req func() error, acceptable Acceptable) error {
	return cb.throttle.doReq(req, nil, acceptable)
}

func (cb *circuitBreaker) DoWithFallback(req func() error, fallback Fallback) error {
	return cb.throttle.doReq(req, fallback, defaultAcceptable)
}

func (cb *circuitBreaker) DoWithFallbackAcceptable(req func() error, fallback Fallback,
	acceptable Acceptable) error {
	return cb.throttle.doReq(req, fallback, acceptable)
}

func (cb *circuitBreaker) Name() string {
	return cb.name
}

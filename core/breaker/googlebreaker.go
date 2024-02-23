package breaker

import (
	"errors"
	"time"
)

var ErrServiceUnavailable = errors.New("circuit breaker is open")

const (
	// 250ms for bucket duration

	//滑动窗口总时间长度，计算失败率和接收率会考虑最近10秒内的请求。
	//如果你将窗口时间增加，熔断器将考虑更长时间的数据，使其对短期波动不那么敏感。相反，减少窗口时间将使熔断器对短期波动更敏感。
	windowDuration = time.Second * 10

	//滑动窗口有多少个桶。每个桶包含了 window 时间段内的请求数据。
	//增加桶的数量将提高统计数据的精度，但同时也会增加内存的使用。减少桶的数量将减少内存使用，但会降低统计数据的精度
	buckets = 40

	//加权系数，计算加权请求数。决定了在计算 dropRatio 时的灵敏度。
	//如果 k 的值较高，熔断器将更宽容地接受请求，即使在一定比例的请求失败时也是如此
	//相反，较低的 k 值会使熔断器在失败率上升时更快地开始拒绝请求，这会让熔断器变得更加敏感。k 的调整可以帮助平衡服务的可用性和对失败的快速响应之间的关系。
	k = 1.5

	//protection 参数定义了在窗口时间内最小的请求总数，该熔断器算法在此之上才开始计算和决策是否需要熔断
	//这是为了避免在请求量很低时由于几个失败请求而触发熔断。如果 protection 的值较低，即使只有很少的请求，熔断器也可能触发。
	//如果 protection 的值较高，那么必须有足够多的请求才会计算失败率和决定是否熔断
	protection = 5
)

type googleBreaker struct {
	k     float64
	stat  *RollingWindow //滑动窗口计数器
	proba *Proba
}

func newGoogleBreaker() *googleBreaker {
	bucketDuration := time.Duration(int64(windowDuration) / int64(buckets))

	st := NewRollingWindow(buckets, bucketDuration)

	return &googleBreaker{
		stat:  st,
		k:     k,
		proba: NewProba(),
	}
}

func (b *googleBreaker) accept() error {
	accepts, total := b.history()
	weightedAccepts := b.k * float64(accepts)

	//https://landing.google.com/sre/sre-book/chapters/handling-overload/#eq2101
	// for better performance, no need to care about negative ratio
	dropRatio := (float64(total-protection) - weightedAccepts) / float64(total+1)
	if dropRatio <= 0 {
		return nil
	}

	if b.proba.TrueOnProba(dropRatio) {
		return ErrServiceUnavailable
	}
	return nil
}

func (b *googleBreaker) allow() (internalPromise, error) {
	if err := b.accept(); err != nil {
		return nil, err
	}

	return googlePromise{
		b: b,
	}, nil
}

func (b *googleBreaker) history() (accepts, total int64) {
	b.stat.Reduce(func(b *Bucket) {
		accepts += int64(b.Sum)
		total += b.Count
	})
	return
}

func (b *googleBreaker) doReq(req func() error, fallback Fallback, acceptable Acceptable) error {
	if err := b.accept(); err != nil {
		b.markFailure()
		if fallback != nil {
			return fallback(err)
		}

		return err
	}

	var success bool
	defer func() {
		// if req() panic, success is false, mark as failure
		if success {
			b.markSuccess()
		} else {
			b.markFailure()
		}
	}()

	err := req()
	if acceptable(err) {
		success = true
	}

	return err
}

func (b *googleBreaker) markSuccess() {
	b.stat.Add(1)
}

func (b *googleBreaker) markFailure() {
	b.stat.Add(0)
}

type googlePromise struct {
	b *googleBreaker
}

func (p googlePromise) Accept() {
	p.b.markSuccess()
}

func (p googlePromise) Reject() {
	p.b.markFailure()
}

package register

import (
	"context"
	"time"
)

// OptionEtcd is etcd registry option.
type OptionEtcd func(o *optionEtcd)

type optionEtcd struct {
	ctx         context.Context
	namespace   string
	ttl         time.Duration
	dialTimeout time.Duration
	maxRetry    int
}

// ContextEtcd with registry context.
func ContextEtcd(ctx context.Context) OptionEtcd {
	return func(o *optionEtcd) { o.ctx = ctx }
}

// Namespace with registry namespace.
func Namespace(ns string) OptionEtcd {
	return func(o *optionEtcd) { o.namespace = ns }
}

// RegisterTTL with register ttl.
func RegisterTTL(ttl time.Duration) OptionEtcd {
	return func(o *optionEtcd) { o.ttl = ttl }
}

func MaxRetry(num int) OptionEtcd {
	return func(o *optionEtcd) { o.maxRetry = num }
}

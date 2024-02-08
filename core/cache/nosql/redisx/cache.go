package redisx

import (
	"context"
	"time"
)

type CacheSession interface {
	TakeCtx(ctx context.Context, value interface{}, key string, queryDB func(val interface{}) error) error
	TakeExpireCtx(ctx context.Context, value interface{}, key string, queryDB func(val interface{}) error) error

	SetCtx(ctx context.Context, key string, val interface{}) error
	SetExpireCtx(ctx context.Context, key string, val interface{}, expire time.Duration) error

	GetCtx(ctx context.Context, key string, val interface{}) error
	DelCtx(ctx context.Context, keys ...string) error
}

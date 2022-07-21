package noSql

import (
	"context"
	"github.com/bhmy-shm/gofks/pkg"
	"time"
)

//专门处理 Redis - string

type StringOperation struct {
	ctx context.Context
}

func StringType() *StringOperation {
	return &StringOperation{ctx: context.Background()}
}

func (this *StringOperation) Set(key string, value interface{}, attrs ...*OperationAttr) *pkg.InterfaceResult {

	//判断是否设置过期时间，UnwrapOr代表如果设置，则 默认30秒 / 或者默认没有过期时间
	exp := OperationAttrs(attrs).Find(EXPR).UnwrapOr(time.Second * 0).(time.Duration)

	//判断是否是nx，
	nx := OperationAttrs(attrs).Find(NX).UnwrapOr(nil)
	if nx != nil {
		return pkg.Result(RedisConn().SetNX(this.ctx, key, value, exp).Result())
	}

	//判断是否是xx，
	xx := OperationAttrs(attrs).Find(XX).UnwrapOr(nil) //如果没有XX 就返回nil
	if xx != nil {
		return pkg.Result(RedisConn().SetXX(this.ctx, key, value, exp).Result())
	}

	return pkg.Result(RedisConn().Set(this.ctx, key, value, exp).Result())
}

func (this *StringOperation) Get(key string) *pkg.InterfaceResult {
	return pkg.Result(RedisConn().Get(this.ctx, key).Result())
}

func (this *StringOperation) MGet(keys ...string) *pkg.InterfaceResult {
	return pkg.Result(RedisConn().MGet(this.ctx, keys...).Result())
}

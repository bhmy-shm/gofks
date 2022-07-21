package noSql

import (
	"errors"
	"regexp"
	"time"
)

/*
	1.使用正则匹配限制传参范围
	2.如果被访问key在，数据库中没有，可以在redis中置入空值，并设置过期时间
	3.ip限流+LUR淘汰
*/

type CachePolicy interface {
	Before(key string) error         //读取缓存前的校验操作
	IfNil(key string, v interface{}) //判断命中数据是否为空
	SetOperation(opt Operation)      //设置redis操作类型
	SetTactics(regexp string)        //设置正则匹配可以访问的key
}

const (
	RegexpInterLength string = `^\d{1,5}$` //匹配m,n位长度的数字
)

//抽象实体

type CrossPolicy struct {
	KeyRegx string        //检查key的正则
	Expire  time.Duration //过期时间
	opt     Operation     //操作类型
}

func NewCrossPolicy(keyRegx string, expire time.Duration) *CrossPolicy {
	return &CrossPolicy{KeyRegx: keyRegx, Expire: expire}
}

func (this *CrossPolicy) Before(key string) error {
	if !regexp.MustCompile(this.KeyRegx).MatchString(key) {
		return errors.New("error cache key")
	}
	return nil
}

func (this *CrossPolicy) IfNil(key string, v interface{}) {
	this.opt.Set(key, v, WithExpire(this.Expire)).Unwrap()
}

//设置策略类型
func (this *CrossPolicy) SetOperation(opt Operation) {
	this.opt = opt
}

//设置正则匹配策略
func (this *CrossPolicy) SetTactics(regexp string) {
	this.KeyRegx = regexp
}

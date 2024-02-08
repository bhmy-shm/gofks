package redisx

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

// 定义一些常量，用于表示黑白名单和空缓存等策略应用
type blackType string

const (
	blacklist blackType = "BlackList-Set"
	whitelist blackType = "WhiteList-Set"
)
const (
	blacklistPrefix = "blacklist:" // 黑名单前缀，后面跟IP地址或用户
	whitelistPrefix = "whitelist:" // 白名单前缀，后面跟IP地址或用户
	cachePrefix     = "cache:"     // 缓存前缀，后面跟请求参数
	cacheEmpty      = "empty"      // 缓存空值标识

	maxAccessTimes  = 10          // 最大访问次数，超过则加入黑名单
	maxAccessPeriod = time.Minute // 最大访问周期，超过则重置访问次数

	cacheTTL = 10 * time.Minute   // 缓存有效期
	blackTTL = 24 * time.Hour     // 黑名单有效期
	whiteTTL = 7 * 24 * time.Hour // 白名单有效期
)

type Checker interface {
	AddToBlackList(string) error //加入黑名单
	AddToWhiteList(string) error //加入白名单

	RemoveWhiteList(string) error //从白名单中删除
	RemoveBlackList(string) error //从黑名单中删除

	CheckWhiteList(string) (bool, error) //判断是否在白名单
	CheckBlackList(string) (bool, error) //判断是否在黑名单

	FindList(string) []string //名单查看,传常量进来
}

type RedisWrapper interface {
	SetExpire(key string, duration time.Duration) bool

	GetAll(key string) map[string]string
	Sets(key string, value map[string]string) error
}

type BlackWhiteUser struct {
	client *redis.Client
	Ctx    context.Context
}

func NewBlackWhite() *BlackWhiteUser {
	return &BlackWhiteUser{
		client: redisClient,
		Ctx:    context.Background(),
	}
}

// SetExpire 设置过期时间
func (this *BlackWhiteUser) SetExpire(key string, duration time.Duration) bool {
	ok, err := this.client.Expire(this.Ctx, key, duration).Result()
	if err != nil {
		return false
	}
	return ok
}

// AddToBlackList 加入黑名单
func (this *BlackWhiteUser) AddToBlackList(user string) error {
	return this.client.SAdd(context.Background(), string(blacklist), user).Err()
}

// AddToWhiteList 加入白名单
func (this *BlackWhiteUser) AddToWhiteList(user string) error {
	return this.client.SAdd(context.Background(), string(whitelist), user).Err()
}

// RemoveWhiteList 更新白名单状态和有效期
func (this *BlackWhiteUser) RemoveWhiteList(user string) error {
	return this.client.SRem(this.Ctx, string(whitelist), user).Err()
}

// RemoveBlackList 更新黑名单状态和有效期
func (this *BlackWhiteUser) RemoveBlackList(user string) error {
	return this.client.SRem(this.Ctx, string(blacklist), user).Err()
}

// FindList 名单查看,传常量进来
func (this *BlackWhiteUser) FindList(types blackType) []string {
	var result []string
	switch types {
	case blacklist:
		result = this.client.SMembers(this.Ctx, string(blacklist)).Val()
	case whitelist:
		result = this.client.SMembers(this.Ctx, string(whitelist)).Val()
	default:
		return nil
	}
	return result
}

// CheckBlackList 用于检查请求IP是否在黑名单中，如果是则返回true，否则返回false
func (this *BlackWhiteUser) CheckBlackList(user string) (bool, error) {
	err := redisClient.SIsMember(context.Background(), string(blacklist), user).Err()
	if err == redis.Nil { // 如果键不存在，则说明不在黑名单中
		return false, nil
	}
	if err != nil { // 如果发生其他错误，则返回错误信息
		return false, err
	}
	return true, nil //说明在黑名单中
}

// CheckWhiteList 用于检查请求IP是否在白名单中，如果是则返回true，否则返回false
func (this *BlackWhiteUser) CheckWhiteList(user string) (bool, error) {
	err := redisClient.SIsMember(context.Background(), string(whitelist), user).Err()
	if err == redis.Nil { // 如果键不存在，则说明不在白名单中
		return false, nil
	}
	if err != nil { // 如果发生其他错误，则返回错误信息
		return false, err
	}
	return true, nil //说明在白名单中
}

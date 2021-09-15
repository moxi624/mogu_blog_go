package common

import (
	"context"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-redis/redis/v8"
	"time"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/10 3:41 下午
 * @version 1.0
 */
var addr, _ = beego.AppConfig.String("redisConn")
var password, _ = beego.AppConfig.String("redisPwd")
var db, _ = beego.AppConfig.Int("db")

var redisClient = redis.NewClient(&redis.Options{
	Addr:     addr,
	Password: password,
	DB:       db, // use default DB
})

var ctx = context.Background()

type redisUtil struct{}

func (*redisUtil) Get(key string) string {
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}
func (*redisUtil) SetEx(key string, value string, expiration int64, unit time.Duration) {
	err := redisClient.SetEX(ctx, key, value, time.Duration(expiration)*unit).Err()
	if err != nil {
		panic(err)
	}
}

func (*redisUtil) Delete(key string) {
	redisClient.Del(ctx, key)
}

func (*redisUtil) Keys(pattern string) []string {
	val, err := redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return []string{}
	}
	return val
}

func (*redisUtil) MultiGet(keys []string) []interface{} {
	val, err := redisClient.MGet(ctx, keys...).Result()
	if err != nil {
		return []interface{}{}
	}
	return val
}

func (*redisUtil) MultiDelete(key []string) {
	redisClient.Del(ctx, key...)
}

func (*redisUtil) Set(key string, value string) {
	err := redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

func (*redisUtil) IncrBy(key string, increment int64) int64 {
	val, err := redisClient.IncrBy(ctx, key, increment).Result()
	if err != nil {
		panic(err)
	}
	return val
}

var RedisUtil = &redisUtil{}

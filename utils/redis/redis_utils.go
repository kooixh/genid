package redis

import (
	"context"
	r "github.com/go-redis/redis/v8"
)

const redisHost = "localhost"
const redisPort = "6379"
const redisPass = ""
const redisDBId = 0
var ctx context.Context
var client *r.Client

func init() {
	ctx = context.Background()
	client = newClient()
}

func newClient() *r.Client {
	return r.NewClient(&r.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPass,
		DB:       redisDBId,
	})
}

func Ping() *r.StatusCmd {
	return client.Ping(ctx)
}

func Set(key string, value interface{}) *r.StatusCmd {
	return client.Set(ctx, key, value, 0)
}

func Get(key string) *r.StringCmd {
	return client.Get(ctx, key)
}

func RPush(key string, value interface{}) *r.IntCmd {
	return client.RPush(ctx, key, value)
}

func LPop(key string) *r.StringCmd {
	return client.LPop(ctx, key)
}

func Length(key string) *r.IntCmd {
	return client.LLen(ctx, key)
}

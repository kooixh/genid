package redis

import (
	"context"
	r "github.com/go-redis/redis/v8"
	"time"
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

func Set(key string, value interface{}, ttl int64) *r.StatusCmd {
	exp := time.Minute * time.Duration(ttl)
	return client.Set(ctx, key, value, exp)
}

func Get(key string) *r.StringCmd {
	return client.Get(ctx, key)
}
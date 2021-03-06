package redis

import (
	"context"
	r "github.com/go-redis/redis/v8"
	c "github.com/kooixh/genid/pkg/constants"
)

var ctx context.Context
var client *r.Client

func init() {
	ctx = context.Background()
	client = newClient()
}

func newClient() *r.Client {
	return r.NewClient(&r.Options{
		Addr:     c.RedisHost + ":" + c.RedisPort,
		Password: c.RedisPass,
		DB:       c.RedisDBId,
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

func Peek(key string) *r.StringSliceCmd {
	return client.LRange(ctx, key, 0, 0)
}

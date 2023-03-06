package rediscounter

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisCounter struct {
	counterKey string
	client     *redis.Client
}

func New(addr string, password string, counterKey string) *RedisCounter {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password,
		DB: 0,
	})
	return &RedisCounter{client: rdb, counterKey: counterKey}
}

func (rc *RedisCounter) Next() (uint64, error) {
	var ctx = context.Background()
	has, err := rc.client.Exists(ctx, rc.counterKey).Result()

	if err != nil {
		return 0, err
	}

	var next int64 = 0
	if has == 1 {
		next, err = rc.client.Incr(ctx, rc.counterKey).Result()
		if err != nil {
			return 0, err
		}
		return uint64(next), nil
	} else {
		_, err = rc.client.Set(ctx, rc.counterKey, 0, 0).Result()
		return 0, err
	}
}

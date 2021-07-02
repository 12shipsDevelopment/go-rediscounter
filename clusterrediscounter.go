package rediscounter

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type ClusterRedisCounter struct {
	counterKey string
	client     *redis.ClusterClient
}

func NewWithCluster(addrs []string, password string, counterKey string) *ClusterRedisCounter {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs,
		Password: password,
	})
	return &ClusterRedisCounter{client: rdb, counterKey: counterKey}
}

func (rc *ClusterRedisCounter) Next() (uint64, error) {
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
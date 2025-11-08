package healthcheck

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisCheck struct {
	rdb *redis.Client
}

func NewRedisCheck(rdb *redis.Client) *RedisCheck {
	return &RedisCheck{
		rdb: rdb,
	}
}

func (r *RedisCheck) Ping(ctx context.Context) error {
	return r.rdb.Ping(ctx).Err()
}

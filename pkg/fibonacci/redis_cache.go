package fibonacci

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type (
	redisCache struct {
		rdb *redis.Client
	}
)

func NewRedisCache(rdb *redis.Client) Cache {
	return &redisCache{rdb}
}

func (r *redisCache) Get(ctx context.Context, n uint8) (string, error) {
	return r.rdb.Get(ctx, strconv.Itoa(int(n))).Result()
}

func (r *redisCache) Save(ctx context.Context, n uint8, val string) error {
	return r.rdb.Set(ctx, strconv.Itoa(int(n)), val, 0).Err()
}

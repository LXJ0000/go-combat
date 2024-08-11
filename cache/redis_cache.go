package _cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	_ "github.com/golang/mock/mockgen/model"
)

type RedisCache struct {
	cmd redis.Cmdable
}

func NewRedisCache(cmd redis.Cmdable) *RedisCache {
	return &RedisCache{
		cmd: cmd,
	}
}

// Get returns the value for the given key.
func (c *RedisCache) Get(ctx context.Context, key string) (any, error) {
	return c.cmd.Get(ctx, key).Result()
}

// Set sets the value for the given key with an optional expiration time.
func (c *RedisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.cmd.Set(ctx, key, value, expiration).Err()
}

// Delete deletes the value for the given key.
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.cmd.Del(ctx, key).Err()
}

// Exists checks if the given key exists in the cache.
func (c *RedisCache) Exists(ctx context.Context, key string) bool {
	return c.cmd.Exists(ctx, key).Val() > 0
}

// LoadAndDelete returns the value for the given key and deletes it from the cache.
func (c *RedisCache) LoadAndDelete(ctx context.Context, key string) (any, error) {
	return c.cmd.GetDel(ctx, key).Result()
}
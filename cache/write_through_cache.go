package _cache

import (
	"context"
	"log/slog"
	"time"
)

type WriteThroughCache struct {
	Cache
	StoreFunc func(ctx context.Context, key string, val any) error
}

func NewWriteThroughCache(store Cache, storeFunc func(ctx context.Context, key string, val any) error) *WriteThroughCache {
	return &WriteThroughCache{
		Cache:     store,
		StoreFunc: storeFunc,
	}
}

// Set writes the value to the cache and the underlying store.
func (c *WriteThroughCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
    if err := c.StoreFunc(ctx, key, val); err != nil {
		return err
	}
	return c.Cache.Set(ctx, key, val, expiration)
}

// SetAsync writes the value to the cache and the underlying store asynchronously.
func (c *WriteThroughCache) SetAsync(ctx context.Context, key string, val any, expiration time.Duration) error {
    if err := c.StoreFunc(ctx, key, val); err != nil {
		return err
	}
	go func ()  {
		if err :=  c.Cache.Set(ctx, key, val, expiration); err != nil {
			slog.Error("write through cache: set data error", slog.String("key", key), slog.Any("val", val), slog.String("error", err.Error()))
		}
	}()
	return nil
}
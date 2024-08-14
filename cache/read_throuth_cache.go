package _cache

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"golang.org/x/sync/singleflight"
)

// ReadThroughCache 必须实现 LoadFunc 以及 Expiration
type ReadThroughCache struct {
	Cache
	LoadFunc   func(ctx context.Context, key string) ([]byte, error)
	Expiration time.Duration
	g          *singleflight.Group
}

func NewReadThroughCache(cache Cache, loadFunc func(ctx context.Context, key string) ([]byte, error), expiration time.Duration) *ReadThroughCache {
	return &ReadThroughCache{
		Cache:      cache,
		LoadFunc:   loadFunc,
		Expiration: expiration,
	}
}

// Get synchronization
func (c *ReadThroughCache) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := c.Cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, errKeyNotFound) {
			value, err = c.LoadFunc(ctx, key)
			if err != nil {
				return nil, err
			}
			if err := c.Cache.Set(ctx, key, value, c.Expiration); err != nil {
				slog.Error("read throuth cache: set data error", slog.String("key", key), slog.String("error", err.Error()))
			}
			return value, nil
		}
		return nil, err
	}
	return value, nil
}

// GetAsync asynchronous
func (c *ReadThroughCache) GetAsync(ctx context.Context, key string) ([]byte, error) {
	value, err := c.Cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, errKeyNotFound) {
			go func() {
				value, err = c.LoadFunc(ctx, key)
				if err != nil {
					slog.Error("read throuth cache: load data error", slog.String("key", key), slog.String("error", err.Error()))
				}
				if err := c.Cache.Set(ctx, key, value, c.Expiration); err != nil {
					slog.Error("read throuth cache: set data error", slog.String("key", key), slog.String("error", err.Error()))
				}
			}()
		}
		return nil, err
	}
	return value, nil
}

// GetAsyncPartial asynchronous
func (c *ReadThroughCache) GetAsyncPartial(ctx context.Context, key string) ([]byte, error) {
	value, err := c.Cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, errKeyNotFound) {
			value, err = c.LoadFunc(ctx, key)
			if err != nil {
				return nil, err
			}
			go func() {
				if err := c.Cache.Set(ctx, key, value, c.Expiration); err != nil {
					slog.Error("read throuth cache: set data error", slog.String("key", key), slog.String("error", err.Error()))
				}
			}()
			return value, nil
		}
		return nil, err
	}
	return value, nil
}

// GetWithSingleflight
func (c *ReadThroughCache) GetWithSingleflight(ctx context.Context, key string) ([]byte, error) {
	value, err := c.Cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, errKeyNotFound) {
			val, err, _ := c.g.Do(key, func() (any, error) {
				value, err := c.LoadFunc(ctx, key)
				if err != nil {
					return nil, err
				}
				if err := c.Cache.Set(ctx, key, value, c.Expiration); err != nil {
					slog.Error("read throuth cache: set data error", slog.String("key", key), slog.String("error", err.Error()))
				}
				return value, nil
			})
			return val.([]byte), err
		}
		return nil, err
	}
	return value, nil
}

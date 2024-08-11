package _cache

import (
	"context"
	"time"

	"golang.org/x/sync/singleflight"
)

type BloomFilterCache struct {
	ReadThroughCache
}

func NewBloomFilterCache(cache Cache, filter BloomFilter,
	loadFunc func(ctx context.Context, key string) (any, error), expiration time.Duration,
) *BloomFilterCache {
	g := &singleflight.Group{}
	return &BloomFilterCache{
		ReadThroughCache: ReadThroughCache{
			Cache: cache,
			LoadFunc: func(ctx context.Context, key string) (any, error) {
				if !filter.Exists(ctx, key) {
					return nil, errKeyNotFound
				}
				return loadFunc(ctx, key)
			},
			Expiration: expiration,
			g: g,
		},
	}
}

type BloomFilter interface {
	Exists(ctx context.Context, key string) bool
}

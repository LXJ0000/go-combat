package _cache

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var (
	errOverCapacity = errors.New("max cnt cache: over capacity")
)

type MaxCntCache struct {
	*LocalCache
	cnt    int32
	maxCnt int32
}

func NewMaxCntCache(maxCnt int32, c *LocalCache) *MaxCntCache {
	cache := &MaxCntCache{
		LocalCache: c,
		cnt:        0,
		maxCnt:     maxCnt,
	}
	originEvict := cache.evict
	cache.evict = func(key string, value any) {
		atomic.AddInt32(&cache.cnt, -1)
		if originEvict != nil {
			originEvict(key, value)
		}
	}
	return cache
}

func (c *MaxCntCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.data[key]
	if !ok {
		if c.cnt+1 > c.maxCnt {
			// 设计淘汰策略
			return errOverCapacity
		}
		c.cnt++
	}
	return c.set(key, value, expiration)
}

// TODO lru
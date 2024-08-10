package _cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	errKeyNotFound = errors.New("key not found")
)

const cleanCount = 1000

type LocalCache struct {
	data  map[string]*item
	mu    sync.RWMutex
	close chan struct{}
}

type item struct {
	value    any
	deadline time.Time
}

// NewLocalCache creates a new LocalCache with the given interval for cleaning up expired items.
func NewLocalCache(interval time.Duration) Cache {
	cache := &LocalCache{
		data: make(map[string]*item),
	}
	// start a goroutine to clean up expired items every interval
	go func() {
		t := time.NewTicker(interval)
		for {
			select {
			case <-cache.close:
				return
			case now := <-t.C:
				cache.mu.Lock()
				i := 0 // count the number of items cleaned up
				for k, v := range cache.data {
					if !v.deadline.IsZero() && now.After(v.deadline) {
						delete(cache.data, k)
					}
					i++
					if i > cleanCount {
						break
					}
				}
				cache.mu.Unlock()
			}
		}
	}()
	return cache
}

// Get returns the value for the given key.
func (c *LocalCache) Get(ctx context.Context, key string) (any, error) {
	c.mu.RLock()
	item, ok := c.data[key]
	c.mu.RUnlock()
	if !ok {
		// return nil, errKeyNotFound
		return nil, fmt.Errorf("local cache: %w, key: %s", errKeyNotFound, key)
	}
	// check if the item is expired
	now := time.Now()
	if !item.deadline.IsZero() && now.After(item.deadline) {
		c.mu.Lock()
		defer c.mu.Unlock()
		item, ok = c.data[key]
		if !ok {
			// return nil, errKeyNotFound
			return nil, fmt.Errorf("local cache: %w, key: %s", errKeyNotFound, key)
		}
		if !item.deadline.IsZero() && now.After(item.deadline) { // double check
			delete(c.data, key)
			return nil, fmt.Errorf("local cache: %w, key: %s", errKeyNotFound, key)
		}
	}
	return item.value, nil
}

// Set sets the value for the given key with an optional expiration time.
// if expiration is zero, the value will not expire.
func (c *LocalCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var deadline time.Time
	if expiration != 0 {
		deadline = time.Now().Add(expiration)
	}
	c.data[key] = &item{
		value:    value,
		deadline: deadline,
	}
	return nil
}

// Delete deletes the value for the given key.
func (c *LocalCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

// Exists checks if the given key exists in the cache.
func (c *LocalCache) Exists(ctx context.Context, key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.data[key]
	return ok, nil
}

// Close closes the cache.
func (c *LocalCache) Close() error {
	select {
	case c.close <- struct{}{}:
		return nil
	default:
		return errors.New("cache is already closed")
	}
}
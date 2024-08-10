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

type LocalCache struct {
	data map[string]any
	mu   sync.RWMutex
}

func NewLocalCache() Cache {
	return &LocalCache{
		data: make(map[string]any),
	}
}

// Get returns the value for the given key.
func (c *LocalCache) Get(ctx context.Context, key string) (any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[key]
	if !ok {
		// return nil, errKeyNotFound
		return nil, fmt.Errorf("local cache: %w, key: %s", errKeyNotFound, key)
	}
	return value, nil
}

// Set sets the value for the given key with an optional expiration time.
func (c *LocalCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
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

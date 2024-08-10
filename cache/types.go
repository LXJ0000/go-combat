package _cache

import (
	"context"
	"time"
)

type Cache interface {
	// Get returns the value for the given key.
	Get(ctx context.Context, key string) (any, error)

	// Set sets the value for the given key with an optional expiration time.
	Set(ctx context.Context, key string, value any, expiration time.Duration) error

	// Delete deletes the value for the given key.
	Delete(ctx context.Context, key string) error

	// Exists checks if the given key exists in the cache.
	Exists(ctx context.Context, key string) bool
}

package _cache

import "sync"

type MaxMemCache struct {
	Cache
	max int64
	use int64
	mu  sync.Mutex
}

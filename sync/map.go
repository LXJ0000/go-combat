package sync

import "sync"

type Map[K comparable, V any] struct {
	values map[K]V
	mu     sync.RWMutex
}

// LoadOrStore 存在则返回，没有则存储
func (m *Map[K, V]) LoadOrStore(k K, v V) (V, bool) {
	m.mu.RLock()
	old, ok := m.values[k]
	m.mu.RUnlock() // defer 则死锁
	if ok {
		return old, true
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	// double check
	if old, ok = m.values[k]; ok {
		return old, true
	}
	m.values[k] = v
	return v, false
}

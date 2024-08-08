package _map

import "sync"

type SyncMap[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		data: make(map[K]V),
	}
}

func (m *SyncMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.data[key]
	return value, ok
}

func (m *SyncMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

func (m *SyncMap[K, V]) LoadOrStore(key K, newValue V) (V, bool) {
	m.mu.RLock()
	oldValue, ok := m.data[key]
	m.mu.RUnlock()
	if ok {
		return oldValue, true
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	oldValue, ok = m.data[key] // dobule check
	if ok {
		return oldValue, true
	}
	m.data[key] = newValue
	return newValue, false
}

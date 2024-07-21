package list

import "sync"

type SyncList[T any] struct {
	List[T]
	mu sync.RWMutex
}

func NewSyncList[T any]() *SyncList[T] {
	return &SyncList[T]{}
}

func (l *SyncList[T]) Cap() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.Len()
}

func (l *SyncList[T]) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.List.Len()
}

func (l *SyncList[T]) Front() T {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.List.Front()
}

func (l *SyncList[T]) Back() T {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.List.Back()
}

func (l *SyncList[T]) PushBack(val T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.List.PushBack(val)
}

func (l *SyncList[T]) PushFront(val T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.List.PushFront(val)
}

func (l *SyncList[T]) PopBack() T {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.List.PopBack()
}

func (l *SyncList[T]) PopFront() T {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.List.PopFront()
}

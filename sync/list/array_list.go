package _list

type ArrayList[T any] struct {
	values []T
}

func (l *ArrayList[T]) Cap() int {
	return cap(l.values)
}

func (l *ArrayList[T]) Len() int {
	return len(l.values)
}
func (l *ArrayList[T]) Front() T {
	return l.values[0]
}

func (l *ArrayList[T]) Back() T {
	return l.values[len(l.values)-1]
}

func (l *ArrayList[T]) PushBack(val T) {
	l.values = append(l.values, val)
}

func (l *ArrayList[T]) PushFront(val T) {
	l.values = append(l.values, val)
}

func (l *ArrayList[T]) PopBack() T {
	if len(l.values) == 0 {
		panic("empty list")
	}
	val := l.values[len(l.values)-1]
	l.values = l.values[:len(l.values)-1]
	return val
}

func (l *ArrayList[T]) PopFront() T {
	if len(l.values) == 0 {
		panic("empty list")
	}
	val := l.values[len(l.values)-1]
	l.values = l.values[:len(l.values)-1]
	return val
}

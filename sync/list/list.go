package list

type List[T any] interface {
	Cap() int
	Len() int
	Front() T
	Back() T
	PushBack(T)
	PushFront(T)
	PopBack() T
	PopFront() T
}

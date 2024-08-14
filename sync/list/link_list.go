package _list

// type List[T any] interface {
// 	Cap() int
// 	Len() int
// 	Front() T
// 	Back() T
// 	PushBack(T)
// 	PushFront(T)
// 	PopBack() T
// 	PopFront() T
// }

type Node[T any] struct {
	value T
	prev  *Node[T]
	next  *Node[T]
}

type LinkList[T any] struct {
	List[T]
	head *Node[T]
	tail *Node[T]
	size int
}

func NewLinkList[T any]() *LinkList[T] {
	head := &Node[T]{}
	tail := &Node[T]{}
	head.next = tail
	tail.prev = head
	return &LinkList[T]{
		head: head,
		tail: tail,
	}
}

func (l *LinkList[T]) Len() int {
	return l.size
}

func (l *LinkList[T]) Empty() bool {
	return l.size == 0
}

func (l *LinkList[T]) Front() T {
	if l.Empty() {
		panic("linklist is empty")
	}
	return l.head.next.value
}

func (l *LinkList[T]) Back() T {
	if l.Empty() {
		panic("linklist is empty")
	}
	return l.tail.prev.value
}

func (l *LinkList[T]) PushBack(value T) {
	node := &Node[T]{value: value}
	node.prev = l.tail.prev
	node.next = l.tail
	l.tail.prev.next = node
	l.tail.prev = node
	l.size++
}

func (l *LinkList[T]) PushFront(value T) {
	node := &Node[T]{value: value}
	node.prev = l.head
	node.next = l.head.next
	l.head.next.prev = node
	l.head.next = node
	l.size++
}

func (l *LinkList[T]) PopBack() T {
	if l.Empty() {
		panic("linklist is empty")
	}
	node := l.tail.prev
	l.Remove(node)
	return node.value
}

func (l *LinkList[T]) PopFront() T {
	if l.Empty() {
		panic("linklist is empty")
	}
	node := l.head.next
	l.Remove(node)
	return node.value
}

func (l *LinkList[T]) Remove(node *Node[T]) {
	if node == nil {
		panic("node can't be nil")
	}
	node.prev.next = node.next
	node.next.prev = node.prev
	l.size--
}

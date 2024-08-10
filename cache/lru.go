package _cache

import "sync"

type Node struct {
	key   string
	value *item
	prev  *Node
	next  *Node
}

type LRU struct {
	mp       map[string]*Node
	head     *Node
	tail     *Node
	capacity int
	mu       sync.RWMutex // TODO Optimize the performance of using the locks
}

func NewLRU(capacity int) *LRU {
	head := &Node{}
	tail := &Node{}
	head.next = tail
	tail.prev = head
	return &LRU{
		mp:       make(map[string]*Node),
		head:     head,
		tail:     tail,
		capacity: capacity,
	}
}

func (lru *LRU) Get(key string) *item {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	node, ok := lru.mp[key]
	if !ok {
		return nil
	}
	lru.moveToHead(node)
	return node.value
}

func (lru *LRU) Put(key string, value *item) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	node, ok := lru.mp[key]
	if ok {
		node.value = value
		lru.moveToHead(node)
		return
	}
	node = &Node{
		key:   key,
		value: value,
	}
	lru.mp[key] = node
	lru.addNode(node)
	if len(lru.mp) > lru.capacity {
		delete(lru.mp, lru.tail.prev.key)
		lru.removeNode(lru.tail.prev)
	}
}

func (lru *LRU) moveToHead(node *Node) {
	lru.removeNode(node)
	lru.addNode(node)
}

func (lru *LRU) addNode(node *Node) {
	node.next = lru.head.next
	lru.head.next.prev = node
	lru.head.next = node
	node.prev = lru.head
}

func (lru *LRU) removeNode(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

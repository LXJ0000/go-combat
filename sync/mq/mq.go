package _mq

import (
	"sync"
)

// 方案1: 每一个消费者订阅的时候创建一个子channel
// 方案2: 轮训所有消费者

type Broker struct {
	mu       sync.RWMutex
	channels []chan Message
	channel  map[string][]chan Message
}

type Message struct {
	Topic   string
	Content string
}

func (b *Broker) Subscribe(cap int) <-chan Message {
	ch := make(chan Message, cap)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.channels = append(b.channels, ch)
	return ch
}

func (b *Broker) SubscribeTopic(topic string, cap int) <-chan Message {
	ch := make(chan Message, cap)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.channel[topic] = append(b.channel[topic], ch)
	return ch
}

func (b *Broker) Send(m Message) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.channels {
		go func() {
			ch <- m
		}()
		// select { // 这种写法 当写满了之后会挂掉
		// case ch <- m:
		// default:
		// 	return fmt.Errorf("channel is full\n")
		// }
	}
	return nil
}

func (b *Broker) Close() error {
	b.mu.Lock()
	channels := b.channels
	b.channels = nil
	b.mu.Unlock()
	for _, ch := range channels { // 避免重复关闭
		close(ch)
	}
	return nil
}

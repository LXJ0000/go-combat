package _mq

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBroker(t *testing.T) {
	b := &Broker{}
	// sender
	go func() {
		for {
			m := Message{Content: "hello world: " + time.Now().String()}
			if err := b.Send(m); err != nil {
				fmt.Printf("send err: %v", err)
				t.Error(err)
				return
			}
			time.Sleep(time.Second)
		}
	}()
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mq := b.Subscribe(0)
			for msg := range mq {
				fmt.Printf("I am %d, receive: %s\n", i, msg.Content)
			}
		}()
	}
	wg.Wait()
}

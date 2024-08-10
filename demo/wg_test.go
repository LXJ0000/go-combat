package demo

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestWG(t *testing.T) {
	wg := sync.WaitGroup{}
	var cnt int64
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&cnt, 1)
		}()
	}
	wg.Wait()
	t.Error(cnt)
}

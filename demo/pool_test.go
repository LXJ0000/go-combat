package demo

import (
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	p := sync.Pool{
		New: func() any {
			return "nil" // 最好别返回nil
		},
	}
	o := p.Get()
	time.Sleep(time.Second)
	p.Put(o)
}

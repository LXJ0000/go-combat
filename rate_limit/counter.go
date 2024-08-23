package ratelimit

import (
	"sync"
	"time"
)

type Counter struct {
	rate  int
	count int
	begin time.Time
	cycle time.Duration
	mu    sync.Mutex
}

func NewCounter(rate int, cycle time.Duration) *Counter {
	return &Counter{
		begin: time.Now(),
		cycle: cycle,
		rate:  rate,
	}
}

func (c *Counter) reset() {
	c.begin = time.Now()
	c.count = 0
}

func (c *Counter) Allow() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	if now.Sub(c.begin) > c.cycle {
		c.reset()
	}

	if c.count >= c.rate {
		return false
	}

	c.count++
	return true
}
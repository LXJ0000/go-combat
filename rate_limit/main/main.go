package main

import (
	"fmt"
	"sync"
	"time"

	ratelimit "github.com/LXJ0000/go-combat/rate_limit"
)

func main() {
	c := ratelimit.NewCounter(3, time.Second)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			if c.Allow() {
				fmt.Println(i, "allowed", time.Now())
			}
		}()
		time.Sleep(time.Millisecond * 200) // 2秒发送10个请求 有6个通过
	}
	wg.Wait()
}

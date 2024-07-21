package demo

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	DoSomething(context.Background())
	ctx := context.Background()
	valCtx := context.WithValue(ctx, "key", map[string]string{})
	valCtx.Value("key")
}

func DoSomething(ctx context.Context) {
	// 和第三方打交道 务必加上 ctx 参数

}

func TestTimeout(t *testing.T) {
	// 一秒内完成业务 MyBusiness
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	done := make(chan struct{})
	go func() {
		MyBusiness()
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("timeout")
	case <-done:
		fmt.Println("done")
	}
}

func MyBusiness() {
	time.Sleep(time.Second * 2)
}

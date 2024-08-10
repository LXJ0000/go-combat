package taskpool

import "context"

type Task func()

type Pool struct {
	task chan Task
	// done atomic.Bool
	done chan struct{}
}

// NewPool gSize goroutine数量 cap 任务队列容量
func NewPool(gSize int, cap int) *Pool {
	pool := &Pool{
		task: make(chan Task, cap),
	}
	for i := 0; i < gSize; i++ {
		go func() {
			// for task := range pool.task {
			// 	// if pool.done.Load() {
			// 	// 	return
			// 	// }
			// 	task()
			// }
			for { // 要是没有 close 则必然会发生泄漏
				select {
				case <-pool.done:
					return
				case task := <-pool.task:
					task()
				}
			}
		}()

	}
	return pool
}

func (p *Pool) Do(t Task) {

}

func (p *Pool) Submit(ctx context.Context, t Task) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.task <- t:
	}
	return nil
}

// Close 释放资源 避免重复调用
func (p *Pool) Close() error {
	// p.done.Store(true)
	close(p.done) // 保证所有 goroutine 都能拿到 <-pool.done 缺陷：重复调用 close 会发生 panic 还是依赖于用户 不管啦 要管就用 Once 限制住
	return nil
}

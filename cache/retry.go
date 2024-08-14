package _cache

import "time"

type RetryStrategy interface {
	Next() (time.Duration, bool)
}

func NewDefaultRetryStrategy(maxRetry int, interval time.Duration) RetryStrategy {
	return &defaultRetryStrategy{
		maxRetry: maxRetry,
		interval: interval,
	}
}

type defaultRetryStrategy struct {
	// 最大重试次数
	maxRetry int
	// 当前重试次数
	retry    int
	// 每次重试间隔时间
	interval time.Duration
}

func (r *defaultRetryStrategy) Next() (time.Duration, bool) {
	if r.retry >= r.maxRetry {
		return 0, false
	}
	r.retry++
	return r.interval, true
}

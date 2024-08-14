package _cache

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	_ "embed"
)

var (
	ErrLockNotFound = errors.New("redis lock: unlock failed with lock not found")
	ErrLockFail     = errors.New("redis lock: lock failed")
	ErrLockRefresh  = errors.New("redis lock: refresh failed")

	//go:embed lua/unlock.lua
	luaUnLock string

	//go:embed lua/refresh_expiration.lua
	luaRefreshExpiration string

	//go:embed lua/lock.lua
	luaLock string
)

type Client struct {
	cmd redis.Cmdable
}

func (c *Client) Lock(ctx context.Context, key string, expiration time.Duration, contextTimeout time.Duration, retry RetryStrategy) (*Lock, error) {
	var ticker *time.Ticker
	value := uuid.New().String() // 唯一标识加锁的人
	for {
		ctxLock, cancel := context.WithTimeout(ctx, contextTimeout)
		res, err := c.cmd.Eval(ctxLock, luaLock, []string{key}, value, expiration.Seconds()).Result()
		cancel()
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				slog.Warn("redis lock: lock failed with context timeout, retrying...")
				continue
			}
			slog.Error("redis lock: lock failed with error", slog.Any("error", err))
			return nil, err
		}
		if res == "OK" {
			return &Lock{
				cmd:        c.cmd,
				key:        key,
				value:      value,
				expiration: expiration,
				done:       make(chan struct{}),
			}, nil
		}
		interval, ok := retry.Next()
		if !ok {
			return nil, ErrLockFail
		}
		if ticker == nil {
			ticker = time.NewTicker(interval)
		} else {
			ticker.Reset(interval)
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (c *Client) TryLock(ctx context.Context, key string, expiration time.Duration) (*Lock, error) {
	value := uuid.New().String() // 唯一标识加锁的人
	ok, err := c.cmd.SetNX(ctx, key, value, expiration).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrLockFail
	}
	return &Lock{
		cmd:        c.cmd,
		key:        key,
		value:      value,
		expiration: expiration,
	}, nil
}

type Lock struct {
	cmd        redis.Cmdable
	key        string
	value      string
	expiration time.Duration
	done       chan struct{}
}

func (l *Lock) UnLock() error {
	defer close(l.done)
	// 以下步骤必须为原子操作 这里采用 lua 脚本实现
	// 1. 检查是否为自己加的锁
	// 2. 解锁
	cnt, err := l.cmd.Eval(context.Background(), luaUnLock, []string{l.key}, l.value).Int64()
	if err != nil {
		return err
	}
	if cnt != 1 {
		return ErrLockNotFound
	}
	return nil
}

func (l *Lock) Refresh(ctx context.Context) error {
	cnt, err := l.cmd.Eval(ctx, luaRefreshExpiration, []string{l.key}, l.value, l.expiration.Seconds()).Int64()
	if err != nil {
		return err
	}
	if cnt != 1 {
		return ErrLockRefresh
	}
	return nil
}

func (l *Lock) AutoRefresh(interval time.Duration, contextTimeout time.Duration) error {
	var lock *Lock // 成功拿到锁
	timeout := make(chan struct{}, 1)
	ticker := time.NewTicker(interval) // 每10秒刷新一次
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
			if err := lock.Refresh(ctx); err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					timeout <- struct{}{}
					cancel()
					continue // 超时之后会进入到 case <-timeout 逻辑
				}
				cancel()
				return err
			}
			cancel()
		case <-l.done:
			return nil // 主动释放锁
		case <-timeout:
			// 超时重试机制...
			ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
			if err := lock.Refresh(ctx); err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					timeout <- struct{}{}
					cancel()
					continue // 超时之后会进入到 case <-timeout 逻辑
				}
				cancel()
				return err
			}
			cancel()
		}
	}
}

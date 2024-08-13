package _cache

import (
	"context"
	"errors"
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
)

type Client struct {
	cmd redis.Cmdable
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
}

func (l *Lock) UnLock() error {
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

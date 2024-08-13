//go:build e2e

package _cache

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRedisCache_e2e_TryLock(t *testing.T) {
	cmd := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "root",
	})

	tcs := []struct {
		name    string
		before  func(t *testing.T) // 准备数据
		after   func(t *testing.T) // 清除数据
		key     string
		wantErr error
		want    *Lock
	}{
		{
			name: "success",
			before: func(t *testing.T) {
			},
			key:     "lock1",
			wantErr: nil,
			want:    &Lock{cmd: cmd, key: "lock1", value: "2"},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				res, err := cmd.GetDel(ctx, "lock1").Result()
				require.NoError(t, err)
				require.NotEmpty(t, res)
			},
		},
		{
			name: "lock by other",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				res, err := cmd.SetNX(ctx, "lock", "1", time.Second*10).Result()
				require.NoError(t, err)
				require.Equal(t, true, res)
			},
			key:     "lock",
			wantErr: ErrLockFail,
			want:    nil,
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				res, err := cmd.GetDel(ctx, "lock").Result()
				require.NoError(t, err)
				require.Equal(t, "1", res)
			},
		},
	}
	for _, tt := range tcs {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(t)
			client := &Client{cmd: cmd}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			lock, err := client.TryLock(ctx, tt.key, time.Minute)
			require.Equal(t, tt.wantErr, err)
			if err == nil {
				require.Equal(t, tt.want.key, lock.key)
				require.NotEmpty(t, lock.value)
				require.NotNil(t, lock.cmd)
			}
			tt.after(t)
		})
	}
}

func TestRedisCache_e2e_UnLock(t *testing.T) {
	cmd := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "root",
	})

	tcs := []struct {
		name    string
		before  func(t *testing.T) // 准备数据
		after   func(t *testing.T) // 清除数据
		wantErr error
		lock    *Lock
	}{
		{
			name:    "lock not exist",
			before:  func(t *testing.T) {},
			after:   func(t *testing.T) {},
			wantErr: ErrLockNotFound, // 锁不存在
			lock: &Lock{
				key:   "lock3",
				cmd:   cmd,
				value: "3",
			},
		},
		{
			name: "lock by other",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				res, err := cmd.SetNX(ctx, "lock4", "4", time.Minute).Result()
				require.NoError(t, err)
				require.True(t, res) // 设置锁成功
			},
			after: func(t *testing.T) {
				_ = cmd.Del(context.Background(), "lock4").Err() // 清除锁
			},
			wantErr: ErrLockNotFound, // 锁不存在
			lock: &Lock{
				key:   "lock4",
				cmd:   cmd,
				value: "555",
			},
		},
		{
			name: "success",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				res, err := cmd.SetNX(ctx, "lock5", "5", time.Minute).Result()
				require.NoError(t, err)
				require.True(t, res) // 设置锁成功
			},
			after: func(t *testing.T) {
				//_ = cmd.Del(context.Background(), "lock5").Err() // 清除锁
			},
			lock: &Lock{
				key:   "lock5",
				cmd:   cmd,
				value: "5",
			},
		},
	}
	for _, tt := range tcs {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(t)
			err := tt.lock.UnLock()
			require.Equal(t, err, tt.wantErr)
			tt.after(t)
		})
	}
}

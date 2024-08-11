// go:build e2e
package _cache

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRedisCache_e2e_Set(t *testing.T) {
	cmd := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "root",
	})
	// cmd.Ping(context.Background()).Result()

	tcs := []struct {
		name    string
		before  func()             // 准备数据
		after   func(t *testing.T) // 清除数据
		key     string
		val     any
		exp     time.Duration
		wantErr error
		want    any
	}{
		{
			name: "success",
			key:  "key",
			val:  "value",
			exp:  time.Minute,
			want: "value",
			after: func(t *testing.T) {
				val, err := cmd.Get(context.Background(), "key").Result()
				require.NoError(t, err)
				require.Equal(t, "value", val)
				err = cmd.Del(context.Background(), "key").Err()
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tcs {
		t.Run(tt.name, func(t *testing.T) {
			// tt.before()
			cache := NewRedisCache(cmd)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			err := cache.Set(ctx, tt.key, tt.val, tt.exp)
			require.NoError(t, err)
			tt.after(t)
		})
	}
}

func TestRedisCache_e2e_Get(t *testing.T) {
	cmd := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "root",
	})
	// cmd.Ping(context.Background()).Result()
	cache := NewRedisCache(cmd)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := cache.Set(ctx, "key", "value", time.Minute)
	require.NoError(t, err)
	val, err := cache.Get(ctx, "key")
	require.NoError(t, err)
	require.Equal(t, "value", val)

}

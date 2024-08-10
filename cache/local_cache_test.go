package _cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalCache_Get(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		cache   func() *LocalCache
		want    any
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "key not found",
			key:  "key1",
			cache: func() *LocalCache {
				return NewLocalCache(time.Second * 10)
			},
			wantErr: fmt.Errorf("local cache: %w, key: %s", errKeyNotFound, "key1"),
		},
		{
			name: "key found",
			key:  "key2",
			cache: func() *LocalCache {
				c := NewLocalCache(time.Second * 10)
				err := c.Set(context.Background(), "key2", "value2", 0)
				require.NoError(t, err)
				return c
			},
			wantErr: nil,
			want:    "value2",
		},
		{
			name: "key expired",
			key:  "key3",
			cache: func() *LocalCache {
				c := NewLocalCache(time.Second*10, func(lc *LocalCache) {})
				err := c.Set(context.Background(), "key3", "value3", time.Second)
				require.NoError(t, err)
				time.Sleep(time.Second * 2)
				return c
			},
			wantErr: fmt.Errorf("local cache: %w, key: %s", errKeyNotFound, "key3"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.cache().Get(context.Background(), tt.key)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, val)
		})
	}
}

func TestLocalCache_Loop(t *testing.T) {
	var cnt int
	c := NewLocalCache(time.Second, WithEvict(func(key string, value any) {
		cnt++
	}))
	err := c.Set(context.Background(), "key1", "value1", time.Second)
	require.NoError(t, err)
	time.Sleep(time.Second * 5)
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.data["key1"]
	require.False(t, ok)
	require.Equal(t, 1, cnt)

}

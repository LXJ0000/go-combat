package _cache

import (
	"context"
	"testing"
	"time"

	"github.com/LXJ0000/go-combat/cache/mocks"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRedisCache_Set(t *testing.T) {
	ts := []struct {
		name       string
		key        string
		value      interface{}
		expiration time.Duration
		mock       func(ctrl *gomock.Controller) redis.Cmdable

		wantErr error
		want    interface{}
	}{
		{
			name:       "success",
			key:        "test",
			value:      "test",
			expiration: time.Second,
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStatusCmd(context.Background())
				status.SetVal("OK")
				cmd.EXPECT().Set(context.Background(), "test", "test", time.Second).Return(status)
				return cmd
			},
		},
		{
			name:       "timeout",
			key:        "test",
			value:      "test",
			expiration: time.Second,
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStatusCmd(context.Background())
				status.SetErr(context.DeadlineExceeded)
				cmd.EXPECT().Set(context.Background(), "test", "test", time.Second).Return(status)
				return cmd
			},
			wantErr: context.DeadlineExceeded,
		},
	}
	for _, tt := range ts {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := &RedisCache{
				cmd: tt.mock(ctrl),
			}
			err := c.Set(context.Background(), tt.key, tt.value, tt.expiration)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisCache_Get(t *testing.T) {
	ts := []struct {
		name string
		key  string
		mock func(ctrl *gomock.Controller) redis.Cmdable

		wantErr error
		want    interface{}
	}{
		{
			name: "success",
			key:  "test",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewStringCmd(context.Background()) // 注意类型匹配
				status.SetVal("test")
				cmd.EXPECT().Get(context.Background(), "test").Return(status)
				return cmd
			},
			want: "test",
		},
	}
	for _, tt := range ts {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := &RedisCache{
				cmd: tt.mock(ctrl),
			}
			val, err := c.Get(context.Background(), tt.key)
			require.Equal(t, tt.wantErr, err)
			if err != nil {
				require.Equal(t, tt.want, val)
			}
		})
	}
}

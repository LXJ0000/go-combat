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

func TestRedisLock_TryLock(t *testing.T) {
	ts := []struct {
		name       string
		key        string
		value      string
		expiration time.Duration
		mock       func(ctrl *gomock.Controller) redis.Cmdable

		wantErr error
		want    *Lock
	}{
		{
			name:       "success",
			key:        "test",
			value:      "test",
			expiration: time.Minute,
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewBoolResult(true, nil)
				cmd.EXPECT().SetNX(context.Background(), "test", gomock.Any(), time.Minute).Return(result)
				return cmd
			},
			want: &Lock{
				key: "test",
			},
		},
		{
			name:       "fail to lock",
			key:        "test",
			value:      "test",
			expiration: time.Minute,
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewBoolResult(false, nil)
				cmd.EXPECT().SetNX(context.Background(), "test", gomock.Any(), time.Minute).Return(result)
				return cmd
			},
			wantErr: ErrLockFail,
		},
		{
			name:       "setnx fail",
			key:        "test",
			value:      "test",
			expiration: time.Minute,
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				result := redis.NewBoolResult(false, context.DeadlineExceeded)
				cmd.EXPECT().SetNX(context.Background(), "test", gomock.Any(), time.Minute).Return(result)
				return cmd
			},
			wantErr: context.DeadlineExceeded,
		},
	}
	for _, tt := range ts {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := &Client{
				cmd: tt.mock(ctrl),
			}
			lock, err := c.TryLock(context.Background(), tt.key, tt.expiration)
			require.Equal(t, tt.wantErr, err)
			if err == nil {
				require.Equal(t, tt.want.key, lock.key)
				require.NotEmpty(t, lock.value)
			}
		})
	}
}

func TestRedisLock_UnLock(t *testing.T) {
	ts := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		key     string
		value   string
		wantErr error
	}{
		{
			name: "lua script eval fail",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewCmd(context.Background())
				status.SetErr(context.DeadlineExceeded)
				cmd.EXPECT().Eval(gomock.Any(), gomock.All(), gomock.Any(), gomock.Any()).
					Return(status)
				return cmd
			},
			key:   "test",
			value: "test",
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "lock not found",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewCmd(context.Background())
				status.SetVal(int64(0))
				cmd.EXPECT().Eval(gomock.Any(), gomock.All(), gomock.Any(), gomock.Any()).
					Return(status)
				return cmd
			},
			key:   "test",
			value: "test",
			wantErr: ErrLockNotFound,
		},
		{
			name: "success",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := mocks.NewMockCmdable(ctrl)
				status := redis.NewCmd(context.Background())
				status.SetVal(int64(1))
				cmd.EXPECT().Eval(gomock.Any(), gomock.All(), gomock.Any(), gomock.Any()).
					Return(status)
				return cmd
			},
			key:   "test",
			value: "test",
		},
	}
	for _, tt := range ts {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			lock := &Lock{
				key:   tt.key,
				value: tt.value,
				cmd:   tt.mock(ctrl),
			}
			err := lock.UnLock()
			require.Equal(t, tt.wantErr, err)
		})
	}
}

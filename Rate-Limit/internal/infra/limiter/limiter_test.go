package limiter

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/bianavic/fullcycle_desafios/internal/infra/limiter/mock"
)

func TestRateLimiter_Allow(t *testing.T) {
	rateLimitIP := 5
	blockTime := 60 * time.Second

	t.Run("success for IP - should return no error if the limit is not exceeded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStrategy := mock.NewMockStorageStrategy(ctrl)

		mockStrategy.EXPECT().Increment(gomock.Any(), "192.168.1.1", blockTime).Return(4, nil)

		tokenConfigs := map[string]TokenConfig{}
		limiter := NewRateLimiter(mockStrategy, rateLimitIP, blockTime, tokenConfigs)

		err := limiter.Allow(context.Background(), "192.168.1.1", "")
		assert.NoError(t, err)
	})

	t.Run("success for token - should return no error if the limit is not exceeded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStrategy := mock.NewMockStorageStrategy(ctrl)

		mockStrategy.EXPECT().Increment(gomock.Any(), "test_token", blockTime).Return(4, nil)

		tokenConfigs := map[string]TokenConfig{
			"test_token": {RateLimit: 100, BlockTime: blockTime},
		}
		limiter := NewRateLimiter(mockStrategy, rateLimitIP, blockTime, tokenConfigs)

		err := limiter.Allow(context.Background(), "", "test_token")
		assert.NoError(t, err)
	})

	t.Run("token not configured - should allow request", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStrategy := mock.NewMockStorageStrategy(ctrl)

		tokenConfigs := map[string]TokenConfig{}
		limiter := NewRateLimiter(mockStrategy, rateLimitIP, blockTime, tokenConfigs)

		err := limiter.Allow(context.Background(), "", "unconfigured_token")
		assert.NoError(t, err)
	})

	t.Run("increment error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStrategy := mock.NewMockStorageStrategy(ctrl)

		mockStrategy.EXPECT().Increment(gomock.Any(), "192.168.1.1", blockTime).Return(0, errors.New("increment error"))

		tokenConfigs := map[string]TokenConfig{}
		limiter := NewRateLimiter(mockStrategy, rateLimitIP, blockTime, tokenConfigs)

		err := limiter.Allow(context.Background(), "192.168.1.1", "")
		assert.Error(t, err)
		assert.Equal(t, "increment error", err.Error())
	})

	t.Run("check rate limit error for token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStrategy := mock.NewMockStorageStrategy(ctrl)

		mockStrategy.EXPECT().Increment(gomock.Any(), "test_token", blockTime).Return(0, errors.New("increment error"))

		tokenConfigs := map[string]TokenConfig{
			"test_token": {RateLimit: 100, BlockTime: blockTime},
		}
		limiter := NewRateLimiter(mockStrategy, rateLimitIP, blockTime, tokenConfigs)

		err := limiter.Allow(context.Background(), "", "test_token")
		assert.Error(t, err)
		assert.Equal(t, "increment error", err.Error())
	})

	t.Run("rate limit exceeded for IP", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStrategy := mock.NewMockStorageStrategy(ctrl)

		mockStrategy.EXPECT().Increment(gomock.Any(), "192.168.1.1", blockTime).Return(6, nil)
		mockStrategy.EXPECT().Set(gomock.Any(), "192.168.1.1", 6, blockTime).Return(nil)

		tokenConfigs := map[string]TokenConfig{}
		limiter := NewRateLimiter(mockStrategy, rateLimitIP, blockTime, tokenConfigs)

		err := limiter.Allow(context.Background(), "192.168.1.1", "")
		assert.Error(t, err)
		assert.Equal(t, ErrRateLimitExceeded, err)
	})

	t.Run("set block time error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockStrategy := mock.NewMockStorageStrategy(ctrl)

		mockStrategy.EXPECT().Increment(gomock.Any(), "192.168.1.1", blockTime).Return(6, nil)
		mockStrategy.EXPECT().Set(gomock.Any(), "192.168.1.1", 6, blockTime).Return(errors.New("set block time error"))

		tokenConfigs := map[string]TokenConfig{}
		limiter := NewRateLimiter(mockStrategy, rateLimitIP, blockTime, tokenConfigs)

		err := limiter.Allow(context.Background(), "192.168.1.1", "")
		assert.Error(t, err)
		assert.Equal(t, "set block time error", err.Error())
	})

	//t.Run("rate limit exceeded for token", func(t *testing.T) {
	//	mockStorage := &MockStorage{
	//		incrementFunc: func(ctx context.Context, key string, expiration time.Duration) (int, error) {
	//			return 101, nil // simulate 101 requests made, exceeding the limit
	//		},
	//	}
	//	tokenConfigs := map[string]TokenConfig{
	//		"test_token": {RateLimit: 100, BlockTime: blockTime},
	//	}
	//	limiter := NewRateLimiter(mockStorage, rateLimitIP, blockTime, tokenConfigs)
	//
	//	err := limiter.Allow(context.Background(), "", "test_token")
	//	assert.Error(t, err)
	//	assert.Equal(t, ErrRateLimitExceeded, err)
	//})
}

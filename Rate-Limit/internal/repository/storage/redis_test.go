package storage

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	mock_storage "github.com/bianavic/fullcycle_desafios/internal/repository/storage/mock"
)

func setupRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return client
}

func TestRedis(t *testing.T) {
	t.Run("should successfully create a redis storage instance", func(t *testing.T) {
		storage, err := NewRedis("localhost:6379", "")
		assert.NoError(t, err)
		assert.NotNil(t, storage)
	})

	t.Run("should fail to get a redis storage instance", func(t *testing.T) {
		storage, err := NewRedis("invalid:6379", "")
		assert.Error(t, err)
		assert.Nil(t, storage)
	})
}

func TestRedis_GetClient(t *testing.T) {
	client := setupRedisClient()
	storage := &Redis{client: client}

	t.Run("GetClient returns valid client", func(t *testing.T) {
		result := storage.GetClient()
		assert.Equal(t, client, result, "Must return the same configured client")
	})
}

func TestRedis_Increment(t *testing.T) {
	client := setupRedisClient()
	storage := &Redis{client: client}
	ctx := context.Background()
	key := "test_key"
	expiration := 5 * time.Second

	t.Run("should successfully increment", func(t *testing.T) {
		client.Del(ctx, key)

		val, err := storage.Increment(ctx, key, expiration)
		assert.NoError(t, err)
		assert.Equal(t, 1, val)

		val, err = storage.Increment(ctx, key, expiration)
		assert.NoError(t, err)
		assert.Equal(t, 2, val)
	})

	t.Run("should return error if increment fails", func(t *testing.T) {
		client.Close()

		_, err := storage.Increment(ctx, key, expiration)

		assert.Error(t, err)
	})

	t.Run("should return error when expire fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := mock_storage.NewMockRedisClient(ctrl)
		storage = &Redis{client: mockClient}

		mockClient.EXPECT().Incr(ctx, key).Return(redis.NewIntResult(1, nil))
		mockClient.EXPECT().Expire(ctx, key, expiration).Return(redis.NewBoolResult(false, fmt.Errorf("failed to set expiration")))

		_, err := storage.Increment(ctx, key, expiration)

		assert.Error(t, err)
		assert.Equal(t, "failed to set expiration", err.Error())
	})
}

func TestRedis_Get(t *testing.T) {
	client := setupRedisClient()
	storage := &Redis{client: client}
	ctx := context.Background()
	key := "test_key"

	t.Run("Get returns zero for non-existent key", func(t *testing.T) {
		client.Del(ctx, key)

		val, err := storage.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, 0, val)
	})

	t.Run("Get returns correct value for existing key", func(t *testing.T) {
		client.Set(ctx, key, 5, 0)

		val, err := storage.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, 5, val)
	})

	t.Run("Get returns error for Redis failure", func(t *testing.T) {
		client.Close()

		_, err := storage.Get(ctx, key)
		assert.Error(t, err)
	})
}

func TestRedis_Set(t *testing.T) {
	client := setupRedisClient()
	storage := &Redis{client: client}
	ctx := context.Background()
	key := "test_key"
	value := 10
	expiration := 5 * time.Second

	t.Run("Set successfully sets value with expiration", func(t *testing.T) {
		err := storage.Set(ctx, key, value, expiration)
		assert.NoError(t, err)

		val, err := client.Get(ctx, key).Int()
		assert.NoError(t, err)
		assert.Equal(t, value, val)

		ttl, err := client.TTL(ctx, key).Result()
		assert.NoError(t, err)
		assert.True(t, ttl > 0)
	})
	//
	//	t.Run("Set returns error for Redis failure", func(t *testing.T) {
	//		client.Close()
	//
	//		err := storage.Set(ctx, key, value, expiration)
	//		assert.Error(t, err)
	//	})
}

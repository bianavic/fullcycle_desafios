package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	t.Run("should successfully create a redis storage instance", func(t *testing.T) {
		storage, err := NewRedis("localhost:6379", "")
		assert.NoError(t, err)
		assert.NotNil(t, storage)
	})
	//
	//t.Run("should successfully ping", func(t *testing.T) {
	//	storage, err := NewRedis("localhost:6379", "")
	//	assert.NoError(t, err)
	//	assert.NotNil(t, storage)
	//
	//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//	defer cancel()
	//	err = storage.GetClient().Ping(ctx).Err()
	//	assert.NoError(t, err)
	//})

	t.Run("should fail to get a redis storage instance", func(t *testing.T) {
		storage, err := NewRedis("invalid:6379", "")
		assert.Error(t, err)
		assert.Nil(t, storage)
	})

	//t.Run("should fail to ping", func(t *testing.T) {
	//	storage, err := NewRedis("localhost:6379", "")
	//	assert.NoError(t, err)
	//	assert.NotNil(t, storage)
	//
	//	storage.GetClient().Close()
	//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//	defer cancel()
	//	err = storage.GetClient().Ping(ctx).Err()
	//	assert.Error(t, err)
	//})
}

//func TestRedis_GetClient(t *testing.T) {
//	client := setupRedisClient()
//	storage := &Redis{client: client}
//
//	t.Run("GetClient returns valid client", func(t *testing.T) {
//		returnedClient := storage.GetClient()
//		assert.NotNil(t, returnedClient)
//		assert.Equal(t, client, returnedClient)
//	})
//
//	t.Run("GetClient after client close", func(t *testing.T) {
//		client.Close()
//		returnedClient := storage.GetClient()
//		assert.NotNil(t, returnedClient)
//	})
//}

//func TestRedis_Increment(t *testing.T) {
//	client := setupRedisClient()
//	storage := &Redis{client: client}
//	ctx := context.Background()
//	key := "test_key"
//	expiration := 5 * time.Second
//
//	client.Del(ctx, key)
//
//	t.Run("should successfully increment", func(t *testing.T) {
//		val, err := storage.Increment(ctx, key, expiration)
//		assert.NoError(t, err)
//		assert.Equal(t, 1, val)
//
//		val, err = storage.Increment(ctx, key, expiration)
//		assert.NoError(t, err)
//		assert.Equal(t, 2, val)
//	})
//
//	t.Run("should set expiration on first increment", func(t *testing.T) {
//		client.Del(ctx, key)
//
//		val, err := storage.Increment(ctx, key, expiration)
//		assert.NoError(t, err)
//		assert.Equal(t, 1, val)
//
//		ttl, err := client.TTL(ctx, key).Result()
//		assert.NoError(t, err)
//		assert.True(t, ttl > 0)
//	})
//
//	t.Run("should return error if setting expiration fails", func(t *testing.T) {
//		client.Del(ctx, key)
//
//		// Simulate failure by using a mock client that returns an error on Expire
//		mockClient := &MockRedisClient{
//			IncrFunc: func(ctx context.Context, key string) *redis.IntCmd {
//				return redis.NewIntResult(1, nil)
//			},
//			ExpireFunc: func(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
//				return redis.NewBoolResult(false, fmt.Errorf("expire error"))
//			},
//		}
//		storage.client = mockClient
//
//		_, err := storage.Increment(ctx, key, expiration)
//		assert.Error(t, err)
//	})
//
//	t.Run("should fail to increment when client is closed", func(t *testing.T) {
//		client.Close()
//
//		_, err := storage.Increment(ctx, key, expiration)
//		assert.Error(t, err)
//	})
//
//	t.Run("should handle redis error", func(t *testing.T) {
//		client.Set(ctx, key, "invalid_value", 0)
//
//		_, err := storage.Increment(ctx, key, expiration)
//		assert.Error(t, err)
//	})
//}

//func TestRedis_Get(t *testing.T) {
//	client := setupRedisClient()
//	storage := &Redis{client: client}
//	ctx := context.Background()
//	key := "test_key"
//
//	t.Run("Get returns zero for non-existent key", func(t *testing.T) {
//		client.Del(ctx, key)
//
//		val, err := storage.Get(ctx, key)
//		assert.NoError(t, err)
//		assert.Equal(t, 0, val)
//	})
//
//	t.Run("Get returns correct value for existing key", func(t *testing.T) {
//		client.Set(ctx, key, 5, 0)
//
//		val, err := storage.Get(ctx, key)
//		assert.NoError(t, err)
//		assert.Equal(t, 5, val)
//	})
//
//	t.Run("Get returns error for Redis failure", func(t *testing.T) {
//		client.Close()
//
//		_, err := storage.Get(ctx, key)
//		assert.Error(t, err)
//	})
//}
//
//func TestRedis_Set(t *testing.T) {
//	client := setupRedisClient()
//	storage := &Redis{client: client}
//	ctx := context.Background()
//	key := "test_key"
//	value := 10
//	expiration := 5 * time.Second
//
//	t.Run("Set successfully sets value with expiration", func(t *testing.T) {
//		err := storage.Set(ctx, key, value, expiration)
//		assert.NoError(t, err)
//
//		val, err := client.Get(ctx, key).Int()
//		assert.NoError(t, err)
//		assert.Equal(t, value, val)
//
//		ttl, err := client.TTL(ctx, key).Result()
//		assert.NoError(t, err)
//		assert.True(t, ttl > 0)
//	})
//
//	t.Run("Set returns error for Redis failure", func(t *testing.T) {
//		client.Close()
//
//		err := storage.Set(ctx, key, value, expiration)
//		assert.Error(t, err)
//	})
//}

//func TestRedis_Increment(t *testing.T) {
//	client := setupRedisClient()
//	storage := &Redis{client: client}
//	ctx := context.Background()
//	key := "test_key"
//	expiration := 5 * time.Second
//
//	client.Del(ctx, key)
//
//	t.Run("should successfully increment", func(t *testing.T) {
//		val, err := storage.Increment(ctx, key, expiration)
//		assert.NoError(t, err)
//		assert.Equal(t, 1, val)
//
//		val, err = storage.Increment(ctx, key, expiration)
//		assert.NoError(t, err)
//		assert.Equal(t, 2, val)
//	})
//
//	t.Run("should set expiration on first increment", func(t *testing.T) {
//		client.Del(ctx, key)
//
//		val, err := storage.Increment(ctx, key, expiration)
//		assert.NoError(t, err)
//		assert.Equal(t, 1, val)
//
//		ttl, err := client.TTL(ctx, key).Result()
//		assert.NoError(t, err)
//		assert.True(t, ttl > 0)
//	})
//
//	t.Run("should return error if setting expiration fails", func(t *testing.T) {
//		client.Del(ctx, key)
//
//		// Simulate failure by using a mock client that returns an error on Expire
//		mockClient := &MockRedisClient{
//			IncrFunc: func(ctx context.Context, key string) *redis.IntCmd {
//				return redis.NewIntResult(1, nil)
//			},
//			ExpireFunc: func(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
//				return redis.NewBoolResult(false, fmt.Errorf("expire error"))
//			},
//		}
//		storage.client = mockClient
//
//		_, err := storage.Increment(ctx, key, expiration)
//		assert.Error(t, err)
//	})
//}

// MockRedisClient is a mock implementation of the Redis client for testing purposes.
type MockRedisClient struct {
	IncrFunc   func(ctx context.Context, key string) *redis.IntCmd
	ExpireFunc func(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	//TODO implement me
	panic("implement me")
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	//TODO implement me
	panic("implement me")
}

func (m *MockRedisClient) TTL(ctx context.Context, key string) *redis.DurationCmd {
	//TODO implement me
	panic("implement me")
}

func (m *MockRedisClient) Close() error {
	//TODO implement me
	panic("implement me")
}

func (m *MockRedisClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	return m.IncrFunc(ctx, key)
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return m.ExpireFunc(ctx, key, expiration)
}

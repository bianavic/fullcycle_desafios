package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// RedisStrategy implements the StorageStrategy interface using Redis.
type RedisStrategy struct {
	Client *redis.Client
}

// NewRedisStrategy creates a new RedisStrategy instance.
func NewRedisStrategy(redisAddr string) *RedisStrategy {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &RedisStrategy{Client: client}
}

// GetClient returns the Redis client instance.
func (r *RedisStrategy) GetClient() *redis.Client {
	return r.Client
}

// Increment increments the value for a key and returns the new value.
func (r *RedisStrategy) Increment(ctx context.Context, key string, expiration time.Duration) (int, error) {
	val, err := r.Client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if val == 1 {
		if err := r.Client.Expire(ctx, key, expiration).Err(); err != nil {
			return 0, err
		}
	}

	return int(val), nil
}

// Get retrieves the value for a key.
func (r *RedisStrategy) Get(ctx context.Context, key string) (int, error) {
	val, err := r.Client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return val, nil
}

// Set sets the value for a key with an expiration time.
func (r *RedisStrategy) Set(ctx context.Context, key string, value int, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

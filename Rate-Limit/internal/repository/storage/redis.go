package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis implements the StorageStrategy interface using Redis.
type Redis struct {
	client RedisClient
}

// NewRedis creates a new Redis instance.
func NewRedis(redisAddr, password string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
		DB:       0,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Redis{client: client}, nil
}

// GetClient returns the Redis client instance.
func (r *Redis) GetClient() RedisClient {
	return r.client
}

// Increment increments the value for a key and returns the new value.
func (r *Redis) Increment(ctx context.Context, key string, expiration time.Duration) (int, error) {
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key %s: %w", key, err)
	}

	// Set expiration if this is the first increment
	if val == 1 {
		if err := r.client.Expire(ctx, key, expiration).Err(); err != nil {
			return 0, fmt.Errorf("failed to set expiration for key %s: %w", key, err)
		}
	}

	return int(val), nil
}

// Get retrieves the value for a key.
func (r *Redis) Get(ctx context.Context, key string) (int, error) {
	val, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, nil
}

// Set sets the value for a key with an expiration time.
func (r *Redis) Set(ctx context.Context, key string, value int, expiration time.Duration) error {
	if err := r.client.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

// FlushAll clears all keys in the Redis database.
func (r *Redis) FlushAll(ctx context.Context) *redis.StatusCmd {
	return r.client.FlushAll(ctx)
}

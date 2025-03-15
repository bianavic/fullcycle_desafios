package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// RedisStorage implements the StorageStrategy interface using Redis.
type RedisStorage struct {
	client *redis.Client
}

// NewRedisStorage creates a new RedisStorage instance.
func NewRedisStorage(addr, password string) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisStorage{client: client}, nil
}

// GetClient returns the Redis client instance.
func (r *RedisStorage) GetClient() *redis.Client {
	return r.client
}

// Increment increments the value for a key and returns the new value.
func (r *RedisStorage) Increment(ctx context.Context, key string, expiration time.Duration) (int, error) {
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	// set expiration if this is the first increment
	if val == 1 {
		if err := r.client.Expire(ctx, key, expiration).Err(); err != nil {
			return 0, err
		}
	}

	return int(val), nil
}

// Get retrieves the requisitions number executed given a specific key.
func (r *RedisStorage) Get(ctx context.Context, key string) (int, error) {
	val, err := r.client.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return val, nil
}

// Set sets the value for a key with an expiration time.
func (r *RedisStorage) Set(ctx context.Context, key string, value int, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

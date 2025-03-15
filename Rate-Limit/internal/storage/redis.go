package storage

import (
	"context"
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

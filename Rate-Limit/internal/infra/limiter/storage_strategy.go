package limiter

import (
	"context"
	"time"
)

//go:generate mockgen -source=storage_strategy.go -destination=./mock/mock_storage_strategy.go - package=mock

// StorageStrategy defines the interface for different storage strategies.
type StorageStrategy interface {
	Increment(ctx context.Context, key string, expiration time.Duration) (int, error)
	Get(ctx context.Context, key string) (int, error)
	Set(ctx context.Context, key string, value int, expiration time.Duration) error
}

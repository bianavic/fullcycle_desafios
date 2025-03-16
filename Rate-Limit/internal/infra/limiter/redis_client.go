package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

//go:generate mockgen -source=redis_client.go -destination=./mock/mock_redis_client.go -package=mock

// RedisClient is an interface that includes the methods used by the Redis struct.
type RedisClient interface {
	Incr(ctx context.Context, key string) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	TTL(ctx context.Context, key string) *redis.DurationCmd
	FlushAll(ctx context.Context) *redis.StatusCmd
	Close() error
}

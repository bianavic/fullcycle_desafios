package usecase

import (
	"context"
	"errors"
	"github.com/bianavic/fullcycle_desafios/internal/db"
	"time"
)

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

// RateLimiter handles rate limiting logic.
type RateLimiter struct {
	storage        db.Storage
	rateLimitIP    int
	rateLimitToken int
	blockTime      time.Duration
}

// NewRateLimiter creates a new RateLimiter instance.
func NewRateLimiter(storage db.Storage, rateLimitIP, rateLimitToken int, blockTime time.Duration) *RateLimiter {
	return &RateLimiter{
		storage:        storage,
		rateLimitIP:    rateLimitIP,
		rateLimitToken: rateLimitToken,
		blockTime:      blockTime,
	}
}

// Allow checks if a request is allowed based on the IP or token.
func (r *RateLimiter) Allow(ctx context.Context, ip, token string) error {
	// check rate limit for token if provided
	if token != "" {
		if err := r.checkRateLimit(ctx, token, r.rateLimitToken); err != nil {
			return err
		}
	}

	// Check rate limit for IP
	return r.checkRateLimit(ctx, ip, r.rateLimitIP)
}

// checkRateLimit checks the rate limit for a given key.
func (r *RateLimiter) checkRateLimit(ctx context.Context, key string, rateLimit int) error {
	count, err := r.storage.Increment(ctx, key, r.blockTime)
	if err != nil {
		return err
	}

	if count > rateLimit {
		return ErrRateLimitExceeded
	}

	return nil
}

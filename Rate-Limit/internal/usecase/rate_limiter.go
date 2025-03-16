package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/bianavic/fullcycle_desafios/internal/repository/storage"
)

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

// RateLimiter handles rate limiting logic.
type RateLimiter struct {
	storage        storage.StorageStrategy
	rateLimitIP    int
	rateLimitToken int
	blockTime      time.Duration
}

// NewRateLimiter creates a new RateLimiter instance.
func NewRateLimiter(storage storage.StorageStrategy, rateLimitIP, rateLimitToken int, blockTime time.Duration) *RateLimiter {
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
			log.Printf("rate limit exceeded for token: %s", token)
			return err
		}
		log.Printf("request allowed for token: %s", token)
		return nil
	}

	// Check rate limit for IP
	return r.checkRateLimit(ctx, ip, r.rateLimitIP)
}

// checkRateLimit checks the rate limit for a given ip.
func (r *RateLimiter) checkRateLimit(ctx context.Context, ip string, rateLimit int) error {
	count, err := r.storage.Increment(ctx, ip, r.blockTime)
	if err != nil {
		log.Printf("Rate limit exceeded for IP: %s", ip)
		return err
	}

	if count > rateLimit {
		return ErrRateLimitExceeded
	}

	log.Printf("request allowed for IP: %s", ip)
	return nil
}

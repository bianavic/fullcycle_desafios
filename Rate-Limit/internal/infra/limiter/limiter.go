package limiter

import (
	"context"
	"errors"
	"log"
	"time"
)

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

// TokenConfig holds the rate limit and expiration time for a specific token.
type TokenConfig struct {
	RateLimit int
	BlockTime time.Duration
}

// RateLimiter handles rate limiting logic.
type RateLimiter struct {
	storage        StorageStrategy
	rateLimitIP    int
	rateLimitToken int
	blockTime      time.Duration
	tokenConfigs   map[string]TokenConfig
}

// NewRateLimiter creates a new RateLimiter instance.
func NewRateLimiter(storage StorageStrategy, rateLimitIP int, blockTime time.Duration, tokenConfigs map[string]TokenConfig) *RateLimiter {
	return &RateLimiter{
		storage:      storage,
		rateLimitIP:  rateLimitIP,
		blockTime:    blockTime,
		tokenConfigs: tokenConfigs,
	}
}

// Allow checks if a request is allowed based on the IP or token.
func (r *RateLimiter) Allow(ctx context.Context, ip, token string) error {
	// check rate limit for token if provided
	if token != "" {
		if config, exists := r.tokenConfigs[token]; exists {
			// use token-specific configuration
			if err := r.checkRateLimit(ctx, token, config.RateLimit, config.BlockTime); err != nil {
				log.Printf("Rate limit exceeded for token: %s", token)
				return err
			}
			log.Printf("Request allowed for token: %s", token)
			return nil
		} else {
			// if the token is not configured, allow the request
			log.Printf("Token %s not configured, allowing request", token)
			return nil
		}
	}

	// check rate limit for IP
	return r.checkRateLimit(ctx, ip, r.rateLimitIP, r.blockTime)
}

// checkRateLimit checks the rate limit for a given ip.
func (r *RateLimiter) checkRateLimit(ctx context.Context, key string, rateLimit int, blockTime time.Duration) error {
	count, err := r.storage.Increment(ctx, key, blockTime)
	if err != nil {
		log.Printf("failed to increment key %s: %v", key, err)
		return err
	}

	if count > rateLimit {
		// set the block time if the rate limit is exceeded
		if err := r.storage.Set(ctx, key, count, blockTime); err != nil {
			log.Printf("failed to set block time for key %s: %v", key, err)
			return err
		}
		return ErrRateLimitExceeded
	}

	log.Printf("request allowed for key: %s", key)
	return nil
}

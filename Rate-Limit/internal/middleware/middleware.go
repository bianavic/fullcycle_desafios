package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/bianavic/fullcycle_desafios/internal/usecase"
)

func RateLimiterMiddleware(limiter *usecase.RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract IP and Token from request
		ip := strings.Split(r.RemoteAddr, ":")[0]
		token := r.Header.Get("API_KEY")
		ctx := context.Background()

		// check rate limit
		err := limiter.Allow(ctx, ip, token)
		if err != nil {
			if errors.Is(err, usecase.ErrRateLimitExceeded) {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			} else {
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

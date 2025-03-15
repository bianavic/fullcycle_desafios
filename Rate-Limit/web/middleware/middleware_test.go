package middleware

import (
	"context"
	"github.com/bianavic/fullcycle_desafios/internal/db"
	"github.com/bianavic/fullcycle_desafios/internal/usecase"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiterMiddleware(t *testing.T) {
	redisStrategy := usecase.NewRedisStrategy("localhost:6379")
	redisClient := redisStrategy.GetClient()

	if err := redisClient.FlushAll(context.Background()).Err(); err != nil {
		t.Fatalf("Failed to clear Redis storage: %v", err)
	}

	redisStorage, err := db.NewRedisStorage("localhost:6379", "")
	if err != nil {
		t.Fatalf("Failed to initialize Redis storage: %v", err)
	}

	rateLimitIP := 1
	rateLimitToken := 1
	blockTime := 1 * time.Minute
	limiter := usecase.NewRateLimiter(redisStorage, rateLimitIP, rateLimitToken, blockTime)

	handler := RateLimiterMiddleware(limiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// test cases
	tests := []struct {
		name       string
		token      string
		expectCode int
	}{
		{"First request (should return status ok)", "test_token", http.StatusOK},
		{"Second request (should be rate limited)", "test_token", http.StatusTooManyRequests},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", server.URL, nil)
			req.Header.Set("API_KEY", tt.token)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectCode {
				t.Errorf("Expected status %d, got %d", tt.expectCode, resp.StatusCode)
			}
		})
	}
}

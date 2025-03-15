package ratelimit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bianavic/fullcycle_desafios/internal/storage"
)

func TestRateLimiterByIP(t *testing.T) {
	// initialize Redis storage
	redisStorage, err := storage.NewRedisStorage("localhost:6379", "")
	if err != nil {
		t.Fatalf("Failed to initialize Redis storage: %v", err)
	}

	// clear Redis storage before test
	if err := redisStorage.GetClient().FlushAll(context.Background()).Err(); err != nil {
		t.Fatalf("Failed to clear Redis storage: %v", err)
	}

	// set rate limit to 5 requests per second for IP
	rateLimitIP := 5
	rateLimitToken := 100
	blockTime := 60 * time.Second
	limiter := NewRateLimiter(redisStorage, rateLimitIP, rateLimitToken, blockTime)

	// create HTTP handler with rate limiter middleware
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
		ip         string
		expectCode int
	}{
		{"First request (should return status ok)", "192.168.1.1", http.StatusOK},
		{"Second request (should return status ok)", "192.168.1.1", http.StatusOK},
		{"Third request (should return status ok)", "192.168.1.1", http.StatusOK},
		{"Fourth request (should return status ok)", "192.168.1.1", http.StatusOK},
		{"Fifth request (should return status ok)", "192.168.1.1", http.StatusOK},
		{"Sixth request (should be rate limited)", "192.168.1.1", http.StatusTooManyRequests},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", server.URL, nil)
			req.RemoteAddr = tt.ip
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

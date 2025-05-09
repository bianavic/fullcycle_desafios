package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	limiter "github.com/bianavic/fullcycle_desafios/internal/infra/limiter"
)

type MockStorage struct{}

func (m *MockStorage) Increment(ctx context.Context, key string, expiration time.Duration) (int, error) {
	return 0, errors.New("storage error")
}

func (m *MockStorage) Get(ctx context.Context, key string) (int, error) {
	return 0, nil
}

func (m *MockStorage) Set(ctx context.Context, key string, value int, expiration time.Duration) error {
	return nil
}

func TestRateLimiterMiddleware(t *testing.T) {
	redisStorage, err := limiter.NewRedis("localhost:6379", "")
	if err != nil {
		t.Fatalf("Failed to initialize Redis storage: %v", err)
	}

	// clear Redis storage before each test
	ctx := context.Background()
	redisStorage.GetClient().FlushAll(ctx)

	rateLimitIP := 1
	blockTime := 1 * time.Minute
	tokenConfigs := map[string]limiter.TokenConfig{
		"test_token": {RateLimit: 1, BlockTime: blockTime},
	}
	l := limiter.NewRateLimiter(redisStorage, rateLimitIP, blockTime, tokenConfigs)

	handler := RateLimiterMiddleware(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	t.Run("Internal server error (should return status 500)", func(t *testing.T) {
		// Use a mock storage that returns an error
		mockStorage := &MockStorage{}
		mockLimiter := limiter.NewRateLimiter(mockStorage, rateLimitIP, blockTime, tokenConfigs)

		// Create a new handler with the mock limiter
		mockHandler := RateLimiterMiddleware(mockLimiter, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}))

		// Create a new test server with the mock handler
		mockServer := httptest.NewServer(mockHandler)
		defer mockServer.Close()

		// Send a request
		req, _ := http.NewRequest("GET", mockServer.URL, nil)
		req.Header.Set("API_KEY", "test_token")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		// Verify the response status code
		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
		}
	})
}

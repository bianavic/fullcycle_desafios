package test

import (
	"context"
	"errors"
	"github.com/bianavic/fullcycle_desafios/internal/usecase"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bianavic/fullcycle_desafios/internal/ratelimit"
	"github.com/bianavic/fullcycle_desafios/internal/repository/storage"
)

type MockStorage struct {
	incrementFunc func(ctx context.Context, key string, expiration time.Duration) (int, error)
}

func (m *MockStorage) Increment(ctx context.Context, key string, expiration time.Duration) (int, error) {
	return m.incrementFunc(ctx, key, expiration)
}

func (m *MockStorage) Get(ctx context.Context, key string) (int, error) {
	return 0, nil
}

func (m *MockStorage) Set(ctx context.Context, key string, value int, expiration time.Duration) error {
	return nil
}

func TestRateLimiter_Allow(t *testing.T) {
	rateLimitIP := 5
	rateLimitToken := 100
	blockTime := 60 * time.Second

	t.Run("success for IP - should return no error if the limit is not exceeded", func(t *testing.T) {
		mockStorage := &MockStorage{
			incrementFunc: func(ctx context.Context, key string, expiration time.Duration) (int, error) {
				return 4, nil // simulate 4 requests made
			},
		}
		limiter := usecase.NewRateLimiter(mockStorage, rateLimitIP, rateLimitToken, blockTime)
		err := limiter.Allow(context.Background(), "192.168.1.1", "")
		assert.NoError(t, err) // Should not return an error if the limit is not exceeded
	})

	t.Run("success for token - should return no error if the limit is not exceeded", func(t *testing.T) {
		mockStorage := &MockStorage{
			incrementFunc: func(ctx context.Context, key string, expiration time.Duration) (int, error) {
				return 99, nil // simulate 99 requests made
			},
		}
		limiter := usecase.NewRateLimiter(mockStorage, rateLimitIP, rateLimitToken, blockTime)
		err := limiter.Allow(context.Background(), "", "test_token")
		assert.NoError(t, err)
	})

	t.Run("rate limit exceeded for IP", func(t *testing.T) {
		mockStorage := &MockStorage{
			incrementFunc: func(ctx context.Context, key string, expiration time.Duration) (int, error) {
				return 6, nil // simulate 6 requests made
			},
		}
		limiter := usecase.NewRateLimiter(mockStorage, rateLimitIP, rateLimitToken, blockTime)
		err := limiter.Allow(context.Background(), "192.168.1.1", "")
		assert.Error(t, err) // should return an error if the limit is exceeded
		assert.Equal(t, usecase.ErrRateLimitExceeded, err)
	})

	t.Run("rate limit exceeded for token", func(t *testing.T) {
		mockStorage := &MockStorage{
			incrementFunc: func(ctx context.Context, key string, expiration time.Duration) (int, error) {
				return 101, nil // simulate 101 requests made
			},
		}
		limiter := usecase.NewRateLimiter(mockStorage, rateLimitIP, rateLimitToken, blockTime)
		err := limiter.Allow(context.Background(), "", "test_token")
		assert.Error(t, err) // should return an error if the limit is exceeded
		assert.Equal(t, usecase.ErrRateLimitExceeded, err)
	})

	t.Run("increment error", func(t *testing.T) {
		mockStorage := &MockStorage{
			incrementFunc: func(ctx context.Context, key string, expiration time.Duration) (int, error) {
				return 0, errors.New("increment error")
			},
		}
		limiter := usecase.NewRateLimiter(mockStorage, rateLimitIP, rateLimitToken, blockTime)
		err := limiter.Allow(context.Background(), "192.168.1.1", "")
		assert.Error(t, err) // Should return an error if there is an increment error
		assert.Equal(t, "increment error", err.Error())
	})
}

func TestRateLimiterByIP(t *testing.T) {
	// initialize Redis storage
	redisStorage, err := storage.NewRedis("localhost:6379", "")
	if err != nil {
		t.Fatalf("Failed to initialize Redis storage: %v", err)
	}
	//
	//// clear Redis storage before test
	//if err := redisStorage.GetClient().FlushAll(context.Background()).Err(); err != nil {
	//	t.Fatalf("Failed to clear Redis storage: %v", err)
	//}

	// set rate limit to 5 requests per second for IP
	rateLimitIP := 5
	rateLimitToken := 100
	blockTime := 60 * time.Second
	limiter := usecase.NewRateLimiter(redisStorage, rateLimitIP, rateLimitToken, blockTime)

	// create HTTP handler with rate limiter middleware
	handler := ratelimit.RateLimiterMiddleware(limiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	t.Run("First request (should return status ok)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.RemoteAddr = "192.168.1.1"
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("Second request (should return status ok)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.RemoteAddr = "192.168.1.1"
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("Third request (should return status ok)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.RemoteAddr = "192.168.1.1"
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("Fourth request (should return status ok)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.RemoteAddr = "192.168.1.1"
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("Fifth request (should return status ok)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.RemoteAddr = "192.168.1.1"
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("Sixth request (should be rate limited)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.RemoteAddr = "192.168.1.1"
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusTooManyRequests {
			t.Errorf("Expected status %d, got %d", http.StatusTooManyRequests, resp.StatusCode)
		}
	})
}

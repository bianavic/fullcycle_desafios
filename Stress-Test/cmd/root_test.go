package cmd

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

type MockHTTPClient struct{}

func (c *MockHTTPClient) Get(url string) (*http.Response, error) {
	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusOK)
	return rec.Result(), nil
}

type MockHTTPClientWithError struct{}

func (c *MockHTTPClientWithError) Get(url string) (*http.Response, error) {
	return nil, errors.New("mock error")
}

func TestWorker(t *testing.T) {

	t.Run("should return 10 successful requests", func(t *testing.T) {
		cfg := Config{
			URL:         "http://test.com",
			Requests:    10,
			Concurrency: 1,
			Client:      &MockHTTPClient{},
		}

		wg := sync.WaitGroup{}
		codes := make(chan int, cfg.Requests)

		wg.Add(1)
		go worker(cfg, &wg, codes, cfg.Requests)

		wg.Wait()
		close(codes)

		successCount := 0
		for code := range codes {
			if code == http.StatusOK {
				successCount++
			}
		}
		if successCount != 10 {
			t.Errorf("Expected 10 successful requests, got %d", successCount)
		}
	})

	t.Run("should return error if request fails", func(t *testing.T) {
		cfg := Config{
			URL:         "http://test.com",
			Requests:    10,
			Concurrency: 1,
			Client:      &MockHTTPClientWithError{},
		}

		wg := sync.WaitGroup{}
		codes := make(chan int, cfg.Requests)

		wg.Add(1)
		go worker(cfg, &wg, codes, cfg.Requests)

		wg.Wait()
		close(codes)

		errorCount := 0
		for code := range codes {
			if code == 0 {
				errorCount++
			}
		}
		if errorCount != 10 {
			t.Errorf("Expected 10 error requests, got %d", errorCount)
		}
	})
}

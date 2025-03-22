package cmd

import (
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

func TestWorker(t *testing.T) {
	cfg := Config{
		URL:         "http://test.com",
		Requests:    10,
		Concurrency: 1,
		Client:      &MockHTTPClient{},
	}

	wg := sync.WaitGroup{}
	codes := make(chan int, cfg.Requests)

	wg.Add(1)
	go worker(cfg, &wg, codes)

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
}

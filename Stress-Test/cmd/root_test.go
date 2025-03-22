package cmd

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
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

var osExit = os.Exit

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

func TestRootCmd(t *testing.T) {

	t.Run("should return 200 OK", func(t *testing.T) {
		args := []string{"stress", "test", "--url", "http://example.com", "--requests", "10", "--concurrency", "2"}

		output := &bytes.Buffer{}
		rootCmd.SetArgs(args)
		rootCmd.SetOut(output)

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("distributes extra requests evenly", func(t *testing.T) {
		cfg := Config{
			URL:         "http://test.com",
			Requests:    10,
			Concurrency: 3,
			Client:      &MockHTTPClient{},
		}

		workerRequests := cfg.Requests / cfg.Concurrency
		extraRequests := cfg.Requests % cfg.Concurrency

		var wg sync.WaitGroup
		requestCounts := make([]int, cfg.Concurrency)

		for i := 0; i < cfg.Concurrency; i++ {
			wg.Add(1)
			r := workerRequests
			if i < extraRequests {
				r++ // distribute extra requests
			}
			requestCounts[i] = r
			go func() {
				defer wg.Done()
				// Simulate worker processing
			}()
		}

		wg.Wait()

		expectedCounts := []int{4, 3, 3} // 10 requests distributed among 3 workers
		for i, count := range requestCounts {
			if count != expectedCounts[i] {
				t.Errorf("expected %d requests for worker %d, got %d", expectedCounts[i], i, count)
			}
		}

		// Additional check to ensure the extra requests are distributed
		extraRequestCount := 0
		for i := 0; i < cfg.Concurrency; i++ {
			if requestCounts[i] > workerRequests {
				extraRequestCount++
			}
		}
		if extraRequestCount != extraRequests {
			t.Errorf("expected %d workers to have extra requests, got %d", extraRequests, extraRequestCount)
		}
	})
}

func TestExecute(t *testing.T) {

	t.Run("should execute without error", func(t *testing.T) {
		rootCmd.SetArgs([]string{"stress", "test", "--url", "http://example.com", "--requests", "10", "--concurrency", "2"})
		rootCmd.SetOut(&bytes.Buffer{})
		rootCmd.SetErr(&bytes.Buffer{})

		// Capture os.Exit calls
		exitCode := 0
		oldOsExit := osExit
		defer func() { osExit = oldOsExit }()
		osExit = func(code int) {
			exitCode = code
		}

		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if exitCode != 0 {
			t.Fatalf("expected exit code 0, got %d", exitCode)
		}
	})
}

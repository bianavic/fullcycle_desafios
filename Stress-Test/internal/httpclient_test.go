package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPClient_Get(t *testing.T) {
	t.Run("Get 200 OK", func(t *testing.T) {
		client := &RealHTTPClient{Client: http.DefaultClient}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, client"))
		}))
		defer ts.Close()

		resp, err := client.Get(ts.URL)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		if resp.Body == nil {
			t.Fatal("expected non-nil body")
		}

		resp.Body.Close()
	})

	t.Run("Get 500 Error", func(t *testing.T) {
		client := &RealHTTPClient{Client: http.DefaultClient}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}))
		defer ts.Close()

		resp, err := client.Get(ts.URL)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
		}

		if resp.Body == nil {
			t.Fatal("expected non-nil body")
		}

		resp.Body.Close()
	})
}

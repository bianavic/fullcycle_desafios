package api

import (
	"context"
	"testing"
)

func TestFetchAddress(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	t.Run("Should successfully create a request when source is ViaCepAPI", func(t *testing.T) {
		_, err := FetchAddress(ctx, "01001000", "ViaCepAPI")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("Should successfully create a request when source is BrasilAPI", func(t *testing.T) {
		_, err := FetchAddress(ctx, "01001000", "BrasilAPI")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("should return error when request with context fails", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel the context immediately

		_, err := FetchAddress(ctx, "01001000", "InvalidAPI")
		if err == nil {
			t.Fatal("Expected error, got none")
		}
	})

	t.Run("Should return error when creating request fails due to invalid URL", func(t *testing.T) {
		_, err := FetchAddress(ctx, "invalid-cep", "InvalidAPI")
		if err == nil {
			t.Fatal("Expected error, got none")
		}
	})

	t.Run("should return error on request failure", func(t *testing.T) {
		_, err := FetchAddress(ctx, "01001000", "InvalidAPI")
		if err == nil {
			t.Fatal("Expected error, got none")
		}
	})

	t.Run("Should return error on decode failure", func(t *testing.T) {
		_, err := FetchAddress(ctx, "00000000", "BrasilAPI")
		if err == nil {
			t.Fatal("Expected error, got none")
		}
	})
}

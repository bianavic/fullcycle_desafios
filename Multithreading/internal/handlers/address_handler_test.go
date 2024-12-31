package handlers

import (
	"context"
	"errors"
	"github.com/bianavic/fullcycle_desafios/infra/api"
	"testing"
	"time"
)

const apiTimeout = 1 * time.Second

type mockAPI struct {
	response interface{}
	err      error
}

func (m *mockAPI) FetchAddress(ctx context.Context, cep, source string) (interface{}, error) {
	return m.response, m.err
}

func successfulViaCepAPIResponse() *api.ViaCepAPIResponse {
	return &api.ViaCepAPIResponse{
		Cep:        "01001000",
		Logradouro: "Praça da Sé",
		Bairro:     "Sé",
		Localidade: "São Paulo",
		Estado:     "SP",
	}
}

func successfulBrasilAPIResponse() *api.BrasilAPIResponse {
	return &api.BrasilAPIResponse{
		Cep:          "01001000",
		Street:       "Praça da Sé",
		Neighborhood: "Sé",
		City:         "São Paulo",
		State:        "SP",
	}
}

func TestFetchAddressHandler_ViaCepAPI_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	t.Run("Should successfully call ViaCepAPI", func(t *testing.T) {

		mock := &mockAPI{response: successfulViaCepAPIResponse()}

		resultChan := make(chan Result)
		go FetchAddressHandler(ctx, "01001000", "ViaCepAPI", mock.FetchAddress, resultChan)

		select {
		case result := <-resultChan:
			if result.Error != nil {
				t.Fatalf("Expected no error, got %v", result.Error)
			}
			if result.Source != "ViaCepAPI" {
				t.Fatalf("Expected source ViaCepAPI, got %s", result.Source)
			}
		case <-time.After(apiTimeout):
			t.Fatal("Expected result, got timeout")
		}
	})

	t.Run("Should successfully call BrasilAPI", func(t *testing.T) {
		mock := &mockAPI{response: successfulBrasilAPIResponse()}

		resultChan := make(chan Result)
		go FetchAddressHandler(ctx, "01001000", "BrasilAPI", mock.FetchAddress, resultChan)

		select {
		case result := <-resultChan:
			if result.Error != nil {
				t.Fatalf("Expected no error, got %v", result.Error)
			}
			if result.Source != "BrasilAPI" {
				t.Fatalf("Expected source BrasilAPI, got %s", result.Source)
			}
		case <-time.After(apiTimeout):
			t.Fatal("Expected result, got timeout")
		}
	})

	t.Run("Should return error when fetching address", func(t *testing.T) {
		mock := &mockAPI{err: errors.New("fetch error")}

		resultChan := make(chan Result)
		go FetchAddressHandler(ctx, "01001000", "BrasilAPI", mock.FetchAddress, resultChan)

		select {
		case result := <-resultChan:
			if result.Error == nil {
				t.Fatal("Expected error, got none")
			}
		case <-time.After(apiTimeout):
			t.Fatal("Expected result, got timeout")
		}
	})
}

func TestPrintResult(t *testing.T) {
	t.Run("Should print ViaCepAPI response correctly", func(t *testing.T) {
		result := Result{
			Source: "ViaCepAPI",
			Data: &api.ViaCepAPIResponse{
				Cep:        "01001000",
				Logradouro: "Praça da Sé",
				Bairro:     "Sé",
				Localidade: "São Paulo",
				Estado:     "SP",
			},
		}
		PrintResult(result)
	})

	t.Run("Should print BrasilAPI response correctly", func(t *testing.T) {
		result := Result{
			Source: "BrasilAPI",
			Data: &api.BrasilAPIResponse{
				Cep:          "01001000",
				Street:       "Praça da Sé",
				Neighborhood: "Sé",
				City:         "São Paulo",
				State:        "SP",
			},
		}
		PrintResult(result)
	})
}

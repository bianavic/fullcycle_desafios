package service

import (
	"errors"
	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFallbackLocationService_GetLocationByCEP(t *testing.T) {

	t.Run("successfully retrieves location from primary service", func(t *testing.T) {
		primary := &mockSuccessfulLocationService{
			Response: &domain.ViaCEPResponse{Cep: "12345678", Localidade: "São Paulo"},
		}
		secondary := &mockFailingLocationService{}

		service := NewFallbackLocationService(primary, secondary)

		resp, err := service.GetLocationByCEP("12345678")
		assert.NoError(t, err)
		assert.Equal(t, "São Paulo", resp.Localidade)
	})

	t.Run("should fallback to secondary service on primary failure", func(t *testing.T) {
		primary := &mockFailingLocationService{}
		secondary := &mockSuccessfulLocationService{
			Response: &domain.ViaCEPResponse{Cep: "12345678", Localidade: "São Paulo"},
		}

		service := NewFallbackLocationService(primary, secondary)

		resp, err := service.GetLocationByCEP("12345678")
		assert.NoError(t, err)
		assert.Equal(t, "São Paulo", resp.Localidade)
	})

	t.Run("returns error when both primary and secondary services fail", func(t *testing.T) {
		primary := &mockFailingLocationService{}
		secondary := &mockFailingLocationService{}

		service := NewFallbackLocationService(primary, secondary)

		resp, err := service.GetLocationByCEP("12345678")
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

type mockFailingLocationService struct{}

func (m *mockFailingLocationService) GetLocationByCEP(cep string) (*domain.ViaCEPResponse, error) {
	return nil, errors.New("mock primary failure")
}

type mockSuccessfulLocationService struct {
	Response *domain.ViaCEPResponse
}

func (m *mockSuccessfulLocationService) GetLocationByCEP(cep string) (*domain.ViaCEPResponse, error) {
	return m.Response, nil
}

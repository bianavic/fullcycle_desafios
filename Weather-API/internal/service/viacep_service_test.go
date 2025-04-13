package service

import (
	"errors"
	"net/http"
	"testing"

	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestViaCEPService_GetLocationByCEPGetLocationByCEP(t *testing.T) {

	t.Run("should successfully get location by CEP", func(t *testing.T) {
		mockHTTPClient := &MockHTTPClient{
			ResponseBody: `{
				"cep":"12345678",
				"localidade":"São Paulo",
				"uf":"SP",
				"bairro":"Centro",
				"logradouro":"Rua Teste"
			}`,
			StatusCode: http.StatusOK,
		}

		service := &ViaCEPService{Client: mockHTTPClient}
		resp, err := service.GetLocationByCEP("12345678")

		assert.NoError(t, err)
		assert.Equal(t, "São Paulo", resp.Localidade)
		assert.Equal(t, "SP", resp.UF)
		assert.Equal(t, "Centro", resp.Bairro)
		assert.Equal(t, "Rua Teste", resp.Logradouro)
	})

	t.Run("should return error when zip code is not found", func(t *testing.T) {
		cep := "00000000"
		mockHTTPClient := &MockHTTPClient{
			ResponseBody: `{"error":"not found"}`,
			StatusCode:   http.StatusNotFound,
		}

		service := NewViaCEPService(mockHTTPClient)
		_, err := service.GetLocationByCEP(cep)

		assert.Error(t, domain.ErrCEPNotFound, err)
	})

	t.Run("should return error when fail to read JSON", func(t *testing.T) {
		mockHTTPClient := &MockHTTPClient{
			ResponseBodyReader: brokenReader{},
			StatusCode:         http.StatusOK,
		}

		service := &ViaCEPService{Client: mockHTTPClient}
		_, err := service.GetLocationByCEP("12345678")

		assert.ErrorIs(t, err, domain.ErrFailedToParseData)

	})

	t.Run("should return error when fail to unmarshalling JSON", func(t *testing.T) {
		cep := "12345678"

		invalidJSON := `{"cep""12345678"}`

		mockHTTPClient := &MockHTTPClient{
			ResponseBody: invalidJSON,
			StatusCode:   http.StatusOK,
		}

		service := &ViaCEPService{Client: mockHTTPClient}
		_, err := service.GetLocationByCEP(cep)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrFailedToParseData)
	})

	t.Run("should  return error when fail to create the HTTP request", func(t *testing.T) {
		invalidCEP := "%%%"

		mockHTTPClient := &MockHTTPClient{
			ResponseBody: "",
			StatusCode:   http.StatusOK,
		}

		service := &ViaCEPService{Client: mockHTTPClient}
		_, err := service.GetLocationByCEP(invalidCEP)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrFailedToCreateRequest)
	})

	t.Run("should return error when fail to send the HTTP request", func(t *testing.T) {
		cep := "12345678"

		mockHTTPClient := &MockHTTPClient{
			Err: errors.New("simulated network error"),
		}

		service := &ViaCEPService{Client: mockHTTPClient}
		_, err := service.GetLocationByCEP(cep)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrFailedToSendRequest)
	})
}

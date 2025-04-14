package service

import (
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestBrasilAPIService_GetLocationByCEPGetLocationByCEP(t *testing.T) {

	t.Run("should successfully get location by CEP", func(t *testing.T) {
		cep := "12345678"
		mockHTTPClient := &MockHTTPClient{
			ResponseBody: `{
				"cep":"12345678",
				"city":"São Paulo",
				"state":"SP",
				"neighborhood":"Centro",
				"street":"Rua Teste"
			}`,
			StatusCode: http.StatusOK,
		}

		service := NewBrasilAPIService(mockHTTPClient)
		response, err := service.GetLocationByCEP(cep)

		assert.NoError(t, err)
		assert.Equal(t, "São Paulo", response.Localidade)
		assert.Equal(t, "SP", response.UF)
		assert.Equal(t, "Centro", response.Bairro)
		assert.Equal(t, "Rua Teste", response.Logradouro)
	})

	t.Run("should return error when zip code is not found", func(t *testing.T) {
		cep := "87654321"
		mockHTTPClient := &MockHTTPClient{
			ResponseBody: `{"error":"can not find zipcode"}`,
			StatusCode:   http.StatusNotFound,
		}

		service := NewBrasilAPIService(mockHTTPClient)
		_, err := service.GetLocationByCEP(cep)

		assert.Error(t, domain.ErrCEPNotFound, err)
	})

	t.Run("should return error when fail to read JSON", func(t *testing.T) {
		cep := "12345678"

		brokenReader := ioutil.NopCloser(brokenReader{})

		mockHTTPClient := &MockHTTPClient{
			ResponseBodyReader: brokenReader,
			StatusCode:         http.StatusOK,
		}

		service := NewBrasilAPIService(mockHTTPClient)
		_, err := service.GetLocationByCEP(cep)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrFailedToParseData)
	})

	t.Run("should return error when fail to unmarshalling JSON", func(t *testing.T) {
		cep := "12345678"

		invalidJSON := `{"cep""12345678"}`

		mockHTTPClient := &MockHTTPClient{
			ResponseBody: invalidJSON,
			StatusCode:   http.StatusOK,
		}

		service := NewBrasilAPIService(mockHTTPClient)
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

		service := NewBrasilAPIService(mockHTTPClient)
		_, err := service.GetLocationByCEP(invalidCEP)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrFailedToCreateRequest)
	})

	t.Run("should return error when fail to send the HTTP request", func(t *testing.T) {
		cep := "12345678"

		mockHTTPClient := &MockHTTPClient{
			Err: errors.New("simulated network error"),
		}

		service := NewBrasilAPIService(mockHTTPClient)
		_, err := service.GetLocationByCEP(cep)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrFailedToSendRequest)
	})
}

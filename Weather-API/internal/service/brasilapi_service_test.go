package service

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct {
	ResponseBody       string
	ResponseBodyReader io.ReadCloser
	StatusCode         int
	Err                error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	body := m.ResponseBodyReader
	if body == nil {
		body = io.NopCloser(strings.NewReader(m.ResponseBody))
	}

	return &http.Response{
		StatusCode: m.StatusCode,
		Body:       body,
		Header:     make(http.Header),
	}, nil
}

func TestGetLocationByCEP(t *testing.T) {

	t.Run("should successfully get location by CEP", func(t *testing.T) {
		cep := "12345678"
		expectedResponse := &domain.ViaCEPResponse{
			Cep:        "12345678",
			Localidade: "São Paulo",
			UF:         "SP",
			Bairro:     "Centro",
			Logradouro: "Rua Teste",
		}

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

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if response.Cep != expectedResponse.Cep {
			t.Errorf("expected %s, got %s", expectedResponse.Cep, response.Cep)
		}
		if response.Localidade != expectedResponse.Localidade {
			t.Errorf("expected %s, got %s", expectedResponse.Localidade, response.Localidade)
		}
		if response.UF != expectedResponse.UF {
			t.Errorf("expected %s, got %s", expectedResponse.UF, response.UF)
		}
		if response.Bairro != expectedResponse.Bairro {
			t.Errorf("expected %s, got %s", expectedResponse.Bairro, response.Bairro)
		}
		if response.Logradouro != expectedResponse.Logradouro {
			t.Errorf("expected %s, got %s", expectedResponse.Logradouro, response.Logradouro)
		}
	})

	t.Run("should return error when zip code is not found", func(t *testing.T) {
		cep := "87654321"
		mockHTTPClient := &MockHTTPClient{
			ResponseBody: `{"error":"not found"}`,
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

type brokenReader struct{}

func (b brokenReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

func (b brokenReader) Close() error {
	return nil
}

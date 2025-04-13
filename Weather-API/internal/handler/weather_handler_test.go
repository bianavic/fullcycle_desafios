package handler

import (
	"encoding/json"
	"errors"
	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bianavic/fullcycle_desafios/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMakeWeatherHandler(t *testing.T) {
	validCEP := "01001000"
	//invalidCEP := "invalid-cep"
	//validCity := "São Paulo"
	validTempC := 25.5

	mockWeatherUsecase := new(mocks.WeatherUseCase)
	handler := MakeWeatherHandler(mockWeatherUsecase)

	t.Run("should return weather data for a valid CEP", func(t *testing.T) {
		mockWeatherUsecase.EXPECT().GetWeatherByCEP(validCEP).
			Return(map[string]float64{"temp_C": validTempC}, nil).
			Once()

		req, err := http.NewRequest("GET", "/weather?cep="+validCEP, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]float64
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, validTempC, response["temp_C"])

		mockWeatherUsecase.AssertExpectations(t)
	})

	t.Run("should return error when CEP is missing", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/weather", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), domain.ErrCEPNotFound.Error())
	})

	t.Run("should return error when WeatherUsecase fails", func(t *testing.T) {
		mockWeatherUsecase.EXPECT().GetWeatherByCEP(validCEP).
			Return(nil, errors.New("unexpected failure")).
			Once()

		req, err := http.NewRequest("GET", "/weather?cep=01001000", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "internal server error")

		mockWeatherUsecase.AssertExpectations(t)
	})

	t.Run("should return error when fail to encode response", func(t *testing.T) {
		mockWeatherUsecase.EXPECT().GetWeatherByCEP(validCEP).
			Return(map[string]float64{"temp_C": validTempC}, nil).
			Once()

		req, err := http.NewRequest("GET", "/weather?cep="+validCEP, nil)
		assert.NoError(t, err)

		faultyWriter := &faultyResponseWriter{header: http.Header{}}

		handler := handler.ServeHTTP
		handler(faultyWriter, req)

		assert.Equal(t, http.StatusInternalServerError, faultyWriter.status)
		mockWeatherUsecase.AssertExpectations(t)
	})
}

type faultyResponseWriter struct {
	header http.Header
	status int
}

func (f *faultyResponseWriter) Header() http.Header {
	return f.header
}

func (f *faultyResponseWriter) WriteHeader(statusCode int) {
	f.status = statusCode
}

func (f *faultyResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("simulated write error")
}

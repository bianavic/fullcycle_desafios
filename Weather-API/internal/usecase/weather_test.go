package usecase

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/bianavic/fullcycle_desafios/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestWeatherUsecase_GetWeatherByCEP(t *testing.T) {
	validCEP := "01001000"
	formattedCEP := "01001000"
	validCity := "São Paulo"
	validTempC := 25.5

	t.Run("should successfully return the right temperature when zip code is valid", func(t *testing.T) {
		mockLocation := new(mocks.LocationService)
		mockWeather := new(mocks.WeatherService)

		mockLocation.EXPECT().GetLocationByCEP("01001000").
			Return(&domain.ViaCEPResponse{
				Localidade: validCity,
			}, nil)

		escapedCity := url.QueryEscape("São Paulo")
		mockWeather.EXPECT().GetWeatherByCity(escapedCity).
			Return(&domain.WeatherAPIResponse{
				Current: struct {
					TempC     float64 `json:"temp_c"`
					TempF     float64 `json:"temp_f"`
					Condition struct {
						Text string `json:"text"`
					} `json:"condition"`
				}{
					TempC: 25.5,
					TempF: 77.9,
					Condition: struct {
						Text string `json:"text"`
					}{
						Text: "Partly cloudy",
					},
				},
			}, nil)

		uc := NewWeatherUsecase(mockLocation, mockWeather, "fake-key")

		result, err := uc.GetWeatherByCEP("01001-000")

		assert.NoError(t, err)
		assert.Equal(t, validTempC, result["temp_C"])

		mockLocation.AssertExpectations(t)
		mockWeather.AssertExpectations(t)
	})

	t.Run("should return error when invalid zip code is provided", func(t *testing.T) {
		mockLocation := new(mocks.LocationService)
		mockWeather := new(mocks.WeatherService)

		uc := NewWeatherUsecase(mockLocation, mockWeather, "fake-key")

		result, err := uc.GetWeatherByCEP("invalid-cep")

		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrInvalidCEP)

		mockLocation.AssertNotCalled(t, "GetLocationByCEP")
		mockWeather.AssertNotCalled(t, "GetWeatherByCity")
	})

	t.Run("should return error when fail to fetch location data", func(t *testing.T) {
		mockLocation := new(mocks.LocationService)
		mockWeather := new(mocks.WeatherService)

		uc := NewWeatherUsecase(mockLocation, mockWeather, "fake-api-key")

		mockLocation.EXPECT().GetLocationByCEP(formattedCEP).
			Return(&domain.ViaCEPResponse{}, domain.ErrFailedLocationData).
			Once()

		result, err := uc.GetWeatherByCEP(validCEP)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrFailedLocationData)

		mockLocation.AssertExpectations(t)
		mockWeather.AssertNotCalled(t, "GetWeatherByCity")
	})

	t.Run("should return error when WeatherService is unavailable", func(t *testing.T) {
		mockLocation := new(mocks.LocationService)
		mockWeather := new(mocks.WeatherService)

		mockLocation.EXPECT().GetLocationByCEP("01001000").
			Return(&domain.ViaCEPResponse{
				Localidade: validCity,
			}, nil)

		escapedCity := url.QueryEscape("São Paulo")

		mockWeather.EXPECT().GetWeatherByCity(escapedCity).
			Return(nil, fmt.Errorf("weather service unavailable"))

		uc := NewWeatherUsecase(mockLocation, mockWeather, "fake-key")

		result, err := uc.GetWeatherByCEP("01001-000")

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrWeatherService)

		mockLocation.AssertExpectations(t)
		mockWeather.AssertExpectations(t)
	})
}

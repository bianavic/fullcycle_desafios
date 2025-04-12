package domain

import "errors"

var (
	ErrInvalidCEP     = errors.New("invalid zipcode")
	ErrCEPNotFound    = errors.New("can not find zipcode")
	ErrWeatherService = errors.New("weather usecase unavailable")
)

func ConvertTemperature(celsius float64) TemperatureResponse {
	return TemperatureResponse{
		TempC: celsius,
		TempF: celsius*1.8 + 32,
		TempK: celsius + 273,
	}
}

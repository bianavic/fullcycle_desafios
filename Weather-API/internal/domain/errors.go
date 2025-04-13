package domain

import "errors"

var (
	ErrInvalidCEP         = errors.New("invalid zipcode")
	ErrCEPNotFound        = errors.New("location not found")
	ErrWeatherService     = errors.New("weather service unavailable")
	ErrAPIKeyMissing      = errors.New("missing WEATHER_API_KEY environment variable")
	ErrFailedLocationData = errors.New("failed to fetch location data")
	ErrFailedWeatherData  = errors.New("failed to fetch weather data")
	ErrFailedToParseData  = errors.New("failed to parse data")
)

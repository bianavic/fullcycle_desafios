package domain

import "errors"

var (
	ErrInvalidCEP            = errors.New("invalid zipcode")
	ErrCEPNotFound           = errors.New("can not find zipcode")
	ErrWeatherService        = errors.New("weather service unavailable")
	ErrAPIKeyMissing         = errors.New("missing WEATHER_API_KEY environment variable")
	ErrFailedLocationData    = errors.New("failed to fetch location data")
	ErrFailedToParseData     = errors.New("failed to parse data")
	ErrFailedToCreateRequest = errors.New("failed to create HTTP request")
	ErrFailedToSendRequest   = errors.New("failed to send HTTP request")
)

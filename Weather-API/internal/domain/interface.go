package domain

//go:generate mockery --name=LocationService --dir=. --output=../mocks --outpkg=mocks --filename=mock_location_service.go
type LocationService interface {
	GetLocationByCEP(cep string) (*ViaCEPResponse, error)
}

//go:generate mockery --name=WeatherService --dir=. --output=../mocks --outpkg=mocks --filename=mock_weather_service.go
type WeatherService interface {
	GetWeatherByCity(city string) (*WeatherAPIResponse, error)
}

package domain

type LocationService interface {
	GetLocationByCEP(cep string) (*LocationResponse, error)
}

type WeatherService interface {
	GetWeatherByCity(city string) (*WeatherAPIResponse, error)
}

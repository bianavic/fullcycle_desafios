package domain

type LocationService interface {
	GetLocationByCEP(cep string) (*ViaCEPResponse, error)
}

type WeatherService interface {
	GetWeatherByCity(city string) (*WeatherAPIResponse, error)
}

package repository

type WeatherRepository interface {
	GetTemperature(city string) (float64, error)
}

package usecase

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/bianavic/fullcycle_desafios/pkg"
)

type WeatherUsecase struct {
	LocationService domain.LocationService
	WeatherService  domain.WeatherService
	APIKey          string
}

func NewWeatherUsecase(locationService domain.LocationService, weatherService domain.WeatherService, apiKey string) *WeatherUsecase {
	return &WeatherUsecase{
		LocationService: locationService,
		WeatherService:  weatherService,
		APIKey:          apiKey,
	}
}

func (uc *WeatherUsecase) GetWeatherByCEP(cep string) (map[string]float64, error) {
	cleanCEP := strings.TrimSpace(strings.ReplaceAll(cep, "-", ""))
	if !regexp.MustCompile(`^\d{8}$`).MatchString(cleanCEP) {
		return nil, domain.ErrInvalidCEP
	}

	location, err := uc.LocationService.GetLocationByCEP(cleanCEP)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedLocationData, err)
	}

	if location.Localidade == "" {
		return nil, domain.ErrCEPNotFound
	}

	city := url.QueryEscape(location.Localidade)
	weatherData, err := uc.WeatherService.GetWeatherByCity(city)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrWeatherService, err)
	}

	tempC := weatherData.Current.TempC
	return pkg.ConvertTemperature(tempC), nil
}

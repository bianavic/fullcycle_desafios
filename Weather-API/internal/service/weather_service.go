package service

import (
	"encoding/json"
	"fmt"
	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"io/ioutil"
	"net/http"
)

type WeatherAPIService struct {
	APIKey string
}

func NewWeatherAPIService(apiKey string) *WeatherAPIService {
	return &WeatherAPIService{APIKey: apiKey}
}

func (s *WeatherAPIService) GetWeatherByCity(city string) (*domain.WeatherAPIResponse, error) {
	weatherURL := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", s.APIKey, city)
	resp, err := http.Get(weatherURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedWeatherData, err)
	}
	defer resp.Body.Close()

	var weatherData domain.WeatherAPIResponse
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedToParseData, err)
	}

	return &weatherData, nil
}

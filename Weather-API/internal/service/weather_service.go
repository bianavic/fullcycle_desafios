package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bianavic/fullcycle_desafios/internal/domain"
)

type WeatherAPIService struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

func NewWeatherAPIService(apiKey string) *WeatherAPIService {
	return &WeatherAPIService{
		APIKey:  apiKey,
		BaseURL: "https://api.weatherapi.com", // URL padrão
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *WeatherAPIService) GetWeatherByCity(city string) (*domain.WeatherAPIResponse, error) {
	if city == "" {
		return nil, fmt.Errorf("city cannot be empty")
	}

	encodedCity := url.QueryEscape(city)
	weatherURL := fmt.Sprintf("%s/v1/current.json?key=%s&q=%s", s.BaseURL, s.APIKey, encodedCity)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(weatherURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedToSendRequest, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read response: %v", domain.ErrFailedToParseData, err)
	}

	if strings.Contains(string(body), "<html>") {
		return nil, fmt.Errorf("%w: weather API returned HTML instead of JSON", domain.ErrWeatherService)
	}

	var weatherData domain.WeatherAPIResponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedToParseData, err)
	}

	return &weatherData, nil
}

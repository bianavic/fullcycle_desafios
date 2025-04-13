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
	APIKey string
}

func NewWeatherAPIService(apiKey string) *WeatherAPIService {
	return &WeatherAPIService{APIKey: apiKey}
}

func (s *WeatherAPIService) GetWeatherByCity(city string) (*domain.WeatherAPIResponse, error) {
	if city == "" {
		return nil, fmt.Errorf("city cannot be empty")
	}

	encodedCity := url.QueryEscape(city)
	weatherURL := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", s.APIKey, encodedCity)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(weatherURL)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %v", domain.ErrFailedWeatherData, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read response: %v", domain.ErrFailedToParseData, err)
	}

	if strings.Contains(string(body), "<html>") {
		return nil, fmt.Errorf("weather API returned HTML error: %s", string(body))
	}

	var weatherData domain.WeatherAPIResponse
	//body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedToParseData, err)
	}

	if weatherData.Current.TempC == 0 {
		return nil, fmt.Errorf("%w: temperature data is zero", domain.ErrFailedWeatherData)
	}

	return &weatherData, nil
}

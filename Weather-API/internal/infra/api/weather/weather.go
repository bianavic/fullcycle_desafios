package weather

import (
	"encoding/json"
	"fmt"
	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/bianavic/fullcycle_desafios/internal/interface/repository"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type WeatherAPIClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewWeatherAPIClient(baseURL string) repository.WeatherRepository {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		panic("WEATHER_API_KEY environment variable is required")
	}

	return &WeatherAPIClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				DisableKeepAlives: true,
			},
		},
	}
}

func (c *WeatherAPIClient) GetTemperature(city string) (float64, error) {
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("%s/v1/current.json?key=%s&q=%s", c.baseURL, c.apiKey, encodedCity)
	log.Printf("Making request to WeatherAPI: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return 0, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("WeatherAPI response status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		var apiError struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&apiError); err == nil && apiError.Error.Message != "" {
			log.Printf("WeatherAPI error: %s", apiError.Error.Message)
			return 0, fmt.Errorf("%w: %s", domain.ErrWeatherService, apiError.Error.Message)
		}
		return 0, fmt.Errorf("%w: status code %d", domain.ErrWeatherService, resp.StatusCode)
	}

	var result struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding response: %v", err)
		return 0, fmt.Errorf("failed to decode weather data: %w", err)
	}

	return result.Current.TempC, nil
}

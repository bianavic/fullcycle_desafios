package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/bianavic/fullcycle_desafios/pkg"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func GetWeatherByCEP(cep string) (map[string]float64, error) {
	cleanCEP := strings.TrimSpace(strings.ReplaceAll(cep, "-", ""))
	if !regexp.MustCompile(`^\d{8}$`).MatchString(cleanCEP) {
		return nil, errors.New("invalid CEP format")
	}

	// ViaCEP
	viaCEPURL := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cleanCEP)
	resp, err := http.Get(viaCEPURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch location data: %w", err)
	}
	defer resp.Body.Close()

	var cepData domain.ViaCEPResponse
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &cepData); err != nil {
		return nil, fmt.Errorf("failed to parse location data: %w", err)
	}

	if cepData.Localidade == "" {
		return nil, errors.New("location not found")
	}

	// WEATHER API
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, errors.New("missing WEATHER_API_KEY environment variable")
	}

	weatherURL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, cepData.Localidade)

	respWeather, err := http.Get(weatherURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer respWeather.Body.Close()

	var weatherData domain.WeatherAPIResponse
	bodyWeather, _ := io.ReadAll(respWeather.Body)
	if err := json.Unmarshal(bodyWeather, &weatherData); err != nil {
		return nil, fmt.Errorf("failed to parse weather data: %w", err)
	}

	// conversion
	tempC := weatherData.Current.TempC
	return pkg.ConvertTemperature(tempC), nil
}

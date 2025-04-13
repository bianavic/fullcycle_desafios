package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/bianavic/fullcycle_desafios/pkg"
)

func GetWeatherByCEP(cep string) (map[string]float64, error) {
	cleanCEP := strings.TrimSpace(strings.ReplaceAll(cep, "-", ""))
	if !regexp.MustCompile(`^\d{8}$`).MatchString(cleanCEP) {
		return nil, domain.ErrInvalidCEP
	}

	// ViaCEP
	viaCEPURL := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cleanCEP)
	resp, err := http.Get(viaCEPURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedLocationData, err)
	}
	defer resp.Body.Close()

	var cepData domain.ViaCEPResponse
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &cepData); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedToParseData, err)
	}

	if cepData.Localidade == "" {
		return nil, domain.ErrCEPNotFound
	}

	// WEATHER API
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, domain.ErrAPIKeyMissing
	}

	city := url.QueryEscape(cepData.Localidade)
	weatherURL := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, city)

	respWeather, err := http.Get(weatherURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedWeatherData, err)
	}
	defer respWeather.Body.Close()

	var weatherData domain.WeatherAPIResponse
	bodyWeather, _ := io.ReadAll(respWeather.Body)
	if err := json.Unmarshal(bodyWeather, &weatherData); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedToParseData, err)
	}

	// conversion
	tempC := weatherData.Current.TempC
	return pkg.ConvertTemperature(tempC), nil
}

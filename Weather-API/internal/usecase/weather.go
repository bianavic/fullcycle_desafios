package usecase

import (
	"errors"
	"fmt"
	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/bianavic/fullcycle_desafios/internal/interface/repository"
	"regexp"
)

type WeatherUseCase struct {
	cepRepo     repository.CEPRepository
	weatherRepo repository.WeatherRepository
}

func NewWeatherUseCase(cepRepo repository.CEPRepository, weatherRepo repository.WeatherRepository) *WeatherUseCase {
	return &WeatherUseCase{
		cepRepo:     cepRepo,
		weatherRepo: weatherRepo,
	}
}

func (uc *WeatherUseCase) GetWeatherByCEP(cep string) (*domain.TemperatureResponse, error) {
	// Validação do CEP
	cleanCEP := regexp.MustCompile(`[^0-9]`).ReplaceAllString(cep, "")
	if len(cleanCEP) != 8 {
		return nil, domain.ErrInvalidCEP
	}

	// Busca localização
	location, err := uc.cepRepo.GetLocation(cleanCEP)
	if err != nil {
		return nil, err
	}

	// Busca temperatura
	tempC, err := uc.weatherRepo.GetTemperature(location.Localidade)
	if err != nil {
		if errors.Is(err, domain.ErrWeatherService) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrWeatherService, err)
	}

	// Converte temperaturas
	result := domain.ConvertTemperature(tempC)
	return &result, nil
}

//func GetWeatherByCEP(cep string) (map[string]float64, error) {
//	cleanCEP := strings.TrimSpace(strings.ReplaceAll(cep, "-", ""))
//	if !regexp.MustCompile(`^\d{8}$`).MatchString(cleanCEP) {
//		return nil, errors.New("invalid CEP format")
//	}
//
//	// ViaCEP
//	viaCEPURL := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cleanCEP)
//	resp, err := http.Get(viaCEPURL)
//	if err != nil {
//		return nil, fmt.Errorf("failed to fetch location data: %w", err)
//	}
//	defer resp.Body.Close()
//
//	var cepData domain.ViaCEPResponse
//	body, _ := ioutil.ReadAll(resp.Body)
//	if err := json.Unmarshal(body, &cepData); err != nil {
//		return nil, fmt.Errorf("failed to parse location data: %w", err)
//	}
//
//	if cepData.Localidade == "" {
//		return nil, errors.New("location not found")
//	}
//
//	// WEATHER API
//	apiKey := os.Getenv("WEATHER_API_KEY")
//	if apiKey == "" {
//		return nil, errors.New("missing WEATHER_API_KEY environment variable")
//	}
//
//	weatherURL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, cepData.Localidade)
//
//	respWeather, err := http.Get(weatherURL)
//	if err != nil {
//		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
//	}
//	defer respWeather.Body.Close()
//
//	var weatherData domain.WeatherAPIResponse
//	bodyWeather, _ := io.ReadAll(respWeather.Body)
//	if err := json.Unmarshal(bodyWeather, &weatherData); err != nil {
//		return nil, fmt.Errorf("failed to parse weather data: %w error (%d): %s", err, respWeather.StatusCode, string(bodyWeather))
//	}
//
//	// conversion
//	tempC := weatherData.Current.TempC
//	return pkg.ConvertTemperature(tempC), nil
//}

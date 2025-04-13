package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/bianavic/fullcycle_desafios/internal/service"
	"github.com/bianavic/fullcycle_desafios/internal/usecase"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	apiKey := getAPIKey()

	locationService := service.NewFallbackLocationService(
		service.NewViaCEPService(),
		service.NewBrasilAPIService(),
	)
	weatherService := service.NewWeatherAPIService(apiKey)
	weatherUsecase := usecase.NewWeatherUsecase(locationService, weatherService, apiKey)

	http.HandleFunc("/weather", makeWeatherHandler(weatherUsecase))

	port := getServerPort()
	fmt.Println("server running on port", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("failed to start server:", err)
	}
}

func loadEnv() error {
	env := os.Getenv("ENV")
	if env == "production" {
		return nil
	}

	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("could not load .env file (dev env): %v", err)
	}

	return nil
}

func getAPIKey() string {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal(domain.ErrAPIKeyMissing)
	}
	log.Printf("server starting with API key: %s", maskAPIKey(apiKey))
	return apiKey
}

func getServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func makeWeatherHandler(weatherUsecase *usecase.WeatherUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cep := r.URL.Query().Get("cep")
		log.Printf("Received request for CEP: %s", cep)

		if cep == "" {
			http.Error(w, domain.ErrMissingCEP.Error(), http.StatusBadRequest)
			return
		}

		result, err := weatherUsecase.GetWeatherByCEP(cep)
		if err != nil {
			log.Printf("Error processing CEP %s: %v", cep, err)
			handleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

func handleError(w http.ResponseWriter, err error) {
	var statusCode int
	var errMsg string

	switch {
	case errors.Is(err, domain.ErrInvalidCEP):
		http.Error(w, domain.ErrInvalidCEP.Error(), http.StatusUnprocessableEntity)
	case errors.Is(err, domain.ErrCEPNotFound):
		http.Error(w, domain.ErrCEPNotFound.Error(), http.StatusNotFound)
	case errors.Is(err, domain.ErrWeatherService):
		http.Error(w, domain.ErrWeatherService.Error(), http.StatusServiceUnavailable)
	case errors.Is(err, domain.ErrFailedLocationData):
		http.Error(w, domain.ErrFailedLocationData.Error(), http.StatusInternalServerError)
	case errors.Is(err, domain.ErrFailedWeatherData):
		http.Error(w, domain.ErrFailedWeatherData.Error(), http.StatusInternalServerError)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	log.Printf("Error occurred: %v", err)
	http.Error(w, errMsg, statusCode)
}

func maskAPIKey(key string) string {
	if len(key) < 8 {
		return "******"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

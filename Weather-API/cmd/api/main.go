package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/bianavic/fullcycle_desafios/internal/service"
	"github.com/bianavic/fullcycle_desafios/internal/usecase"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Print("warning: .env file not found - using system environment variables")
		}
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}
	log.Printf("server starting with API key: %s", maskAPIKey(apiKey))

	locationService := service.NewViaCEPService()
	weatherService := service.NewWeatherAPIService(apiKey)
	weatherUsecase := usecase.NewWeatherUsecase(locationService, weatherService, apiKey)

	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		cep := r.URL.Query().Get("cep")
		if cep == "" {
			http.Error(w, "missing 'cep' parameter", http.StatusBadRequest)
			return
		}

		result, err := weatherUsecase.GetWeatherByCEP(cep)
		if err != nil {
			handleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("server running on port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("failed to start server:", err)
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

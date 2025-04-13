package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/bianavic/fullcycle_desafios/internal/domain"
)

func MakeWeatherHandler(weatherUsecase domain.WeatherUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cep := r.URL.Query().Get("cep")
		log.Printf("Received request for CEP: %s", cep)

		if cep == "" {
			http.Error(w, domain.ErrCEPNotFound.Error(), http.StatusBadRequest)
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
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	log.Printf("Error occurred: %v", err)
	http.Error(w, errMsg, statusCode)
}

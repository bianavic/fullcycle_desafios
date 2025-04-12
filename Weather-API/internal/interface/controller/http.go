package controller

import (
	"encoding/json"
	"errors"
	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/bianavic/fullcycle_desafios/internal/usecase"
	"net/http"
)

type WeatherController struct {
	useCase *usecase.WeatherUseCase
}

func NewWeatherController(useCase *usecase.WeatherUseCase) *WeatherController {
	return &WeatherController{useCase: useCase}
}

func (c *WeatherController) GetWeather(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		http.Error(w, "CEP parameter is required", http.StatusBadRequest)
		return
	}

	temp, err := c.useCase.GetWeatherByCEP(cep)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCEP):
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		case errors.Is(err, domain.ErrCEPNotFound):
			http.Error(w, "can not find zipcode", http.StatusNotFound)
		case errors.Is(err, domain.ErrWeatherService):
			http.Error(w, "weather service unavailable", http.StatusInternalServerError)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(temp)
}

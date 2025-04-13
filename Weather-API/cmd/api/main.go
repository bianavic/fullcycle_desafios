package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bianavic/fullcycle_desafios/internal/usecase"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to the Weather API!"))
	})

	http.HandleFunc("/weather", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	server := &http.Server{
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("server running on port 8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("error starting server:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		http.Error(w, "missing 'cep' parameter", http.StatusBadRequest)
		return
	}

	result, err := usecase.GetWeatherByCEP(cep)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCEP):
			http.Error(w, domain.ErrInvalidCEP.Error(), http.StatusUnprocessableEntity)
		case errors.Is(err, domain.ErrCEPNotFound):
			http.Error(w, domain.ErrCEPNotFound.Error(), http.StatusNotFound)
		case errors.Is(err, domain.ErrWeatherService):
			http.Error(w, domain.ErrWeatherService.Error(), http.StatusServiceUnavailable)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func maskAPIKey(key string) string {
	if len(key) < 8 {
		return "******"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

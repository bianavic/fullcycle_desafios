package main

import (
	"log"
	"net/http"

	"github.com/bianavic/fullcycle_desafios/Weather-API-Otel/config"
	"github.com/bianavic/fullcycle_desafios/Weather-API-Otel/handlers"
	"github.com/bianavic/fullcycle_desafios/Weather-API-Otel/middleware"
)

func main() {

	cfg := config.Load()

	mux := http.NewServeMux()

	// Health Check
	mux.HandleFunc("/health", handlers.HealthCheck)

	// CEP Handler with middleware
	cepHandler := http.HandlerFunc(handlers.CEPHandler)
	mux.Handle("/cep", middleware.ValidateCEP(cepHandler))

	log.Printf("Starting server on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

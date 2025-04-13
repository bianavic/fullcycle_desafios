package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/bianavic/fullcycle_desafios/internal/handler"
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

	http.HandleFunc("/weather", handler.MakeWeatherHandler(weatherUsecase))

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

func maskAPIKey(key string) string {
	if len(key) < 8 {
		return "******"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

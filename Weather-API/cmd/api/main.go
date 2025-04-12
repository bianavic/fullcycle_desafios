package main

import (
	"github.com/bianavic/fullcycle_desafios/internal/infra/api/weather"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bianavic/fullcycle_desafios/internal/infra/api/viacep"
	"github.com/bianavic/fullcycle_desafios/internal/interface/controller"
	"github.com/bianavic/fullcycle_desafios/internal/usecase"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Warning: Error loading .env file - using system environment variables")
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	log.Print("Starting server...")

	// Inicializa os repositórios concretos
	cepRepo := viacep.NewViaCEPClient("https://viacep.com.br")
	weatherRepo := weather.NewWeatherAPIClient("http://api.weatherapi.com")

	// Cria o use case com as dependências injetadas
	weatherUseCase := usecase.NewWeatherUseCase(cepRepo, weatherRepo)

	// Cria o controller
	weatherController := controller.NewWeatherController(weatherUseCase)

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Weather API!"))
	})
	router.HandleFunc("/weather", weatherController.GetWeather)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))

	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: router,
		// Timeouts para evitar conexões pendentes
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server starting on port %s with API key: %s", port, maskAPIKey(apiKey))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
	}
}

// maskAPIKey mascara a maior parte da API key para logs
func maskAPIKey(key string) string {
	if len(key) < 8 {
		return "******"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

//func main() {
//	log.Print("sta\t//err := godotenv.Load()\n\t//if err != nil {\n\t//\tlog.Fatal(\"Error loading .env file\")\n\t//}rting server...")
//
//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		w.Write([]byte("Welcome to the Weather API!"))
//	})
//
//	http.HandleFunc("/weather", handler)
//
//	port := os.Getenv("PORT")
//	if port == "" {
//		port = "8080"
//		log.Printf("defaulting to port %s", port)
//	}
//
//	server := &http.Server{
//		Addr:         ":8080",
//		WriteTimeout: 15 * time.Second,
//		ReadTimeout:  15 * time.Second,
//	}
//
//	fmt.Println("Server running on port 8080")
//	if err := server.ListenAndServe(); err != nil {
//		fmt.Println("Error starting server:", err)
//	}
//}

//func handler(w http.ResponseWriter, r *http.Request) {
//	cep := r.URL.Query().Get("cep")
//	if cep == "" {
//		http.Error(w, "Missing 'cep' parameter", http.StatusBadRequest)
//		return
//	}
//
//	result, err := usecase.GetWeatherByCEP(cep)
//	if err != nil {
//		if err.Error() == "invalid zip code" {
//			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
//		} else if err.Error() == "zip code not found" {
//			http.Error(w, err.Error(), http.StatusNotFound)
//		} else {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//		}
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(result)
//}

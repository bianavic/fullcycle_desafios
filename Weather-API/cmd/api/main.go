package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bianavic/fullcycle_desafios/internal/service"
)

func main() {
	log.Print("sta\t//err := godotenv.Load()\n\t//if err != nil {\n\t//\tlog.Fatal(\"Error loading .env file\")\n\t//}rting server...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Weather API!"))
	})

	http.HandleFunc("/weather", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	server := &http.Server{
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server running on port 8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		http.Error(w, "Missing 'cep' parameter", http.StatusBadRequest)
		return
	}

	result, err := service.GetWeatherByCEP(cep)
	if err != nil {
		if err.Error() == "invalid zip code" {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		} else if err.Error() == "zip code not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

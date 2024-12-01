package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

const (
	apiURL                  = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	serverPort              = ":8080"
	apiTimeout              = 200 * time.Millisecond
	dbFile                  = "exchange_rates.db"
	dbTimeout               = 10 * time.Millisecond                                          // Timeout for the database operation (10ms)
	dbPath                  = "file:/app/data/db/exchange_rates.db?cache=shared&mode=memory" // sqlite db file
	createExchangeRateTable = `CREATE TABLE IF NOT EXISTS exchange_rates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
)

// CurrencyRateRequest structure for API response
type CurrencyRateRequest struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

type CurrencyRateResponse struct {
	Bid string `json:"bid"`
}

func main() {
	db, err := connectToDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/cotacao", func(writer http.ResponseWriter, request *http.Request) {
		ExchangeRateHandler(writer, request, db)
	})
	fmt.Printf("Server running on http://localhost%s/cotacao\n", serverPort)
	if err := http.ListenAndServe(serverPort, nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

// startServer starts an HTTP server that serves exchange rates
//func startServer() {}

func ExchangeRateHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("starting request")
	defer log.Println("request finalized")

	rate, err := fetchExchangeRate()
	if err != nil {
		http.Error(w, "failed to fetch exchange rate", http.StatusInternalServerError)
		fmt.Printf("error fetching exchange rate: %v\n", err)
		return
	}

	resp := CurrencyRateResponse{
		Bid: rate.USDBRL.Bid,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to decode JSON", http.StatusInternalServerError)
		log.Printf("failed to decode JSON: %v", err)
		return
	}

	if err := storeExchangeRate(db, rate); err != nil {
		http.Error(w, "failed to sava data to d", http.StatusInternalServerError)
		log.Printf("Erro ao salvar cotação no banco de dados: %v", err)
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Store the rate in the database
	if err := storeExchangeRate(db, rate); err != nil {
		log.Printf("Error storing exchange rate: %v\n", err)
	}
}

// fetchExchangeRate fetches the exchange rate from the API
func fetchExchangeRate() (CurrencyRateRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return CurrencyRateRequest{}, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{
		Timeout: 3 * time.Minute, // timeout for the client
	}
	resp, err := client.Do(req)
	if err != nil {
		return CurrencyRateRequest{}, fmt.Errorf("failed to fetch exchange rate: %w", err)
	}
	defer resp.Body.Close()

	var rates CurrencyRateRequest
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return CurrencyRateRequest{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return rates, nil
}

func storeExchangeRate(db *sql.DB, rate CurrencyRateRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `INSERT INTO exchange_rates (bid) VALUES (?);`

	_, err := db.ExecContext(ctx, query, rate.USDBRL.Bid)
	if err != nil {
		return fmt.Errorf("failed to insert exchange rate: %w", err)
	}

	return err
}

func connectToDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}
	defer db.Close()

	// Create the exchange_rates table if it doesn't exist
	_, err = db.ExecContext(context.Background(), createExchangeRateTable)
	if err != nil {
		log.Fatalf("Error creating table in SQLite: %v", err)
	}
	return db, err
}

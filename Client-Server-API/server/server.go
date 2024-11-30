package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	apiURL     = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	apiTimeout = 200 * time.Millisecond
	dbFile     = "exchange_rates.db"
)

var db *sql.DB
var mu sync.Mutex // handle concurrent database access

// CurrencyRates structure for API response
type CurrencyRates struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func ExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting request")
	defer log.Println("request finalized")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rate, err := fetchExchangeRate()
	if err != nil {
		http.Error(w, "failed to fetch exchange rate", http.StatusInternalServerError)
		fmt.Printf("error fetching exchange rate: %v\n", err)
		return
	}

	// Store the rate in the database
	if err := storeExchangeRate(rate); err != nil {
		log.Printf("Error storing exchange rate: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"bid": rate})
}

// fetchExchangeRate fetches the exchange rate from the API
func fetchExchangeRate() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{
		Timeout: 3 * time.Minute, // timeout for the client
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch exchange rate: %w", err)
	}
	defer resp.Body.Close()

	var rates CurrencyRates
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return rates.USDBRL.Bid, nil
}

func storeExchangeRate(rate string) error {
	mu.Lock()
	defer mu.Unlock()

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open SQLite database: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO exchange_rates (bid) VALUES (?)", rate)
	if err != nil {
		return fmt.Errorf("failed to insert exchange rate: %w", err)
	}

	return nil
}

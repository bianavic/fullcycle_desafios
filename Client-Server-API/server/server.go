package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	apiURL     = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	apiTimeout = 200 * time.Millisecond
)

type CurrencyRate struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type CurrencyRates struct {
	USDBRL CurrencyRate `json:"usdbrl"`
}

func ExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	moedas, err := fetchExchangeRate()
	if err != nil {
		http.Error(w, "failed to fetch exchange rate", http.StatusInternalServerError)
		fmt.Printf("error fetching exchange rate: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"bid": moedas})
}

func fetchExchangeRate() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch exchange rate: %w", err)
	}
	defer resp.Body.Close()

	var c CurrencyRates
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return c.USDBRL.Bid, nil
}

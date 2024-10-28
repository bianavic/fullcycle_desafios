package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiURL     = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	serverPort = ":8080"
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

func main() {
	http.HandleFunc("/cotacao", ExchangeRateHandler)
	fmt.Printf("Server running on http://localhost%s/cotacao\n", serverPort)
	if err := http.ListenAndServe(serverPort, nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

func ExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	moedas, err := FetchExchangeRate()
	if err != nil {
		http.Error(w, "Failed to fetch exchange rate", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"bid": moedas})
}

func FetchExchangeRate() (string, error) {
	resp, err := http.Get(apiURL)
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

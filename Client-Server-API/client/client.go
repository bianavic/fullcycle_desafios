package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	serverURL = "http://localhost:8080/cotacao"
	timeout   = 300 * time.Millisecond
)

type BidResponse struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rate, err := getExchangeRate(ctx)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	if err := saveToFile(rate); err != nil {
		fmt.Printf("error saving to file: %v\n", err)
		return
	}

	fmt.Println("exchange rate saved to cotacao.txt ")
}

func getExchangeRate(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch exchange rate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}

	var bidResp BidResponse
	if err := json.NewDecoder(resp.Body).Decode(&bidResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return bidResp.Bid, nil
}

func saveToFile(rate string) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dolar: %s", rate))
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
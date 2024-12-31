package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	brasilApiURL = "https://brasilapi.com.br/api/cep/v1/"
	viaCepApiURL = "https://viacep.com.br/ws/"
	apiTimeout   = 1 * time.Second
)

type BrasilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type ViaCepAPIResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Estado     string `json:"estado"`
}

type Result struct {
	Source string
	Data   interface{}
	Error  error
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Add the zip code to the command: go run main.go <zip code>")
		return
	}
	cep := os.Args[1]

	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	resultChan := make(chan Result)

	go fetchAddress(ctx, cep, brasilApiURL+cep, resultChan, "BrasilAPI")
	go fetchAddress(ctx, cep, viaCepApiURL+cep+"/json/", resultChan, "ViaCepAPI")

	select {
	case result := <-resultChan:
		if result.Error != nil {
			fmt.Printf("Error fetching address: %v\n", result.Error)
		}
		printResult(result)
	case <-time.After(apiTimeout):
		fmt.Println("Timeout")
	}
}

func fetchAddress(ctx context.Context, cep, url string, resultChan chan Result, source string) {

	//if source == "ViaCepAPI" {
	//	time.Sleep(time.Millisecond * 2000) // Simulate delay for ViaCepAPI
	//}
	//if source == "BrasilAPI" {
	//	time.Sleep(time.Millisecond * 2000) // Simulate delay for BrasilAPI
	//}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resultChan <- Result{Source: source, Error: fmt.Errorf("error creating request: %v", err)}
		return
	}
	client := &http.Client{Timeout: apiTimeout}
	resp, err := client.Do(req)
	if err != nil {
		resultChan <- Result{Source: source, Error: fmt.Errorf("error fetching address: %v", err)}
		return
	}
	defer resp.Body.Close()

	var address interface{}
	if source == "BrasilAPI" {
		address = &BrasilAPIResponse{}
	} else {
		address = &ViaCepAPIResponse{}
	}

	if err := json.NewDecoder(resp.Body).Decode(address); err != nil {
		resultChan <- Result{Source: source, Error: fmt.Errorf("error decoding response: %v", err)}
		return
	}

	resultChan <- Result{Source: source, Data: address}
}

func printResult(result Result) {
	switch data := result.Data.(type) {
	case *ViaCepAPIResponse:
		fmt.Printf("Received from viacep: source:%s - %s, %s - %s, %s, %s\n", result.Source, data.Logradouro, data.Bairro, data.Localidade, data.Estado, data.Cep)
	case *BrasilAPIResponse:
		fmt.Printf("Received from brasilapi: source:%s - %s, %s - %s, %s, %s\n", result.Source, data.Street, data.Neighborhood, data.City, data.State, data.Cep)
	}
}

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

func FetchAddress(ctx context.Context, cep, source string) (interface{}, error) {
	url := getAPIURL(cep, source)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	client := &http.Client{Timeout: apiTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching address: %v", err)
	}
	defer resp.Body.Close()

	address, err := decodeResponse(resp, source)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func getAPIURL(cep, source string) string {
	if source == "BrasilAPI" {
		return brasilApiURL + cep
	}
	return viaCepApiURL + cep + "/json/"
}

func decodeResponse(resp *http.Response, source string) (interface{}, error) {
	var address interface{}
	if source == "BrasilAPI" {
		//time.Sleep(time.Second * 4) // Simulate delay for BrasilAPI
		address = &BrasilAPIResponse{}
	} else {
		//time.Sleep(time.Second * 4) // Simulate delay for ViaCepAPI
		address = &ViaCepAPIResponse{}
	}

	if err := json.NewDecoder(resp.Body).Decode(address); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return address, nil
}

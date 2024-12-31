package main

import (
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

type AddressResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state,omitempty" json:"estado,omitempty"`
	City         string `json:"city,omitempty" json:"localidade,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty" json:"bairro,omitempty"`
	Street       string `json:"street,omitempty" json:"logradouro,omitempty"`
}

func main() {
	cep := "16300025"
	address, err := fetchAddressFromBrasilAPI(cep)
	if err != nil {
		fmt.Printf("Erro ao buscar endereço: %v\n", err)
		return
	}

	fmt.Printf("Endereço da API: %s, %s - %s, %s, %s\n", address.Street, address.Neighborhood, address.City, address.State, address.Cep)
}

func fetchAddressFromBrasilAPI(cep string) (*AddressResponse, error) {
	req, err := http.NewRequest(http.MethodGet, brasilApiURL+cep, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: apiTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var address AddressResponse
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return nil, err
	}

	return &address, nil
}

func fetchAddressFromViaCepAPI(cep string) (*AddressResponse, error) {
	req, err := http.NewRequest(http.MethodGet, viaCepApiURL+cep+"/json", nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: apiTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var address AddressResponse
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return nil, err
	}

	return &address, nil
}

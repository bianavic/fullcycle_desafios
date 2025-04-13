package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bianavic/fullcycle_desafios/internal/domain"
)

type ViaCEPService struct {
	Client domain.HTTPClient
}

func NewViaCEPService(client domain.HTTPClient) *ViaCEPService {
	return &ViaCEPService{Client: client}
}

func (s *ViaCEPService) GetLocationByCEP(cep string) (*domain.ViaCEPResponse, error) {
	viaCEPURL := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequest(http.MethodGet, viaCEPURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedLocationData, err)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedLocationData, err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedToParseData, err)
	}

	if erro, ok := raw["erro"].(bool); ok && erro {
		return nil, domain.ErrCEPNotFound
	}

	var viaCEPData struct {
		Localidade string `json:"localidade"`
		UF         string `json:"uf"`
		CEP        string `json:"cep"`
		Bairro     string `json:"bairro"`
		Logradouro string `json:"logradouro"`
	}

	if err := json.Unmarshal(body, &viaCEPData); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedToParseData, err)
	}

	if viaCEPData.Localidade == "" {
		return nil, domain.ErrCEPNotFound
	}

	return &domain.ViaCEPResponse{
		Cep:        viaCEPData.CEP,
		Localidade: viaCEPData.Localidade,
		UF:         viaCEPData.UF,
		Bairro:     viaCEPData.Bairro,
		Logradouro: viaCEPData.Logradouro,
	}, nil
}

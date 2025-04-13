package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bianavic/fullcycle_desafios/internal/domain"
)

type ViaCEPService struct{}

func NewViaCEPService() *ViaCEPService {
	return &ViaCEPService{}
}

func (s *ViaCEPService) GetLocationByCEP(cep string) (*domain.LocationResponse, error) {
	viaCEPURL := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(viaCEPURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedLocationData, err)
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
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

	return &domain.LocationResponse{
		City:     viaCEPData.Localidade,
		State:    viaCEPData.UF,
		CEP:      viaCEPData.CEP,
		District: viaCEPData.Bairro,
		Street:   viaCEPData.Logradouro,
		Service:  "ViaCEP",
	}, nil
}

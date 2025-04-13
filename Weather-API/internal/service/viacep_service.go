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

func (s *ViaCEPService) GetLocationByCEP(cep string) (*domain.ViaCEPResponse, error) {
	viaCEPURL := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(viaCEPURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedLocationData, err)
	}
	defer resp.Body.Close()

	var cepData domain.ViaCEPResponse
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &cepData); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedToParseData, err)
	}

	return &cepData, nil
}

package service

import (
	"encoding/json"
	"fmt"
	"io"
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
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedToCreateRequest, err)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedToSendRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, domain.ErrCEPNotFound
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedToParseData, err)
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

	return &domain.ViaCEPResponse{
		Cep:        viaCEPData.CEP,
		Localidade: viaCEPData.Localidade,
		UF:         viaCEPData.UF,
		Bairro:     viaCEPData.Bairro,
		Logradouro: viaCEPData.Logradouro,
	}, nil
}

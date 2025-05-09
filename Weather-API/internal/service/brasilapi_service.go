package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bianavic/fullcycle_desafios/internal/domain"
)

type BrasilAPIService struct {
	Client domain.HTTPClient
}

func NewBrasilAPIService(client domain.HTTPClient) *BrasilAPIService {
	return &BrasilAPIService{
		Client: client,
	}
}

func (s *BrasilAPIService) GetLocationByCEP(cep string) (*domain.ViaCEPResponse, error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	var data struct {
		Cep          string `json:"cep"`
		City         string `json:"city"`
		State        string `json:"state"`
		Neighborhood string `json:"neighborhood"`
		Street       string `json:"street"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedToParseData, err)
	}

	return &domain.ViaCEPResponse{
		Cep:        data.Cep,
		Localidade: data.City,
		UF:         data.State,
		Bairro:     data.Neighborhood,
		Logradouro: data.Street,
	}, nil
}

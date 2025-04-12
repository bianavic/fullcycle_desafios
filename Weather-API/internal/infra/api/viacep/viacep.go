package viacep

import (
	"encoding/json"
	"fmt"
	"github.com/bianavic/fullcycle_desafios/internal/domain"
	"github.com/bianavic/fullcycle_desafios/internal/interface/repository"
	"net/http"
	"time"
)

type ViaCEPClient struct {
	baseURL string
	client  *http.Client
}

func NewViaCEPClient(baseURL string) repository.CEPRepository {
	return &ViaCEPClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *ViaCEPClient) GetLocation(cep string) (*domain.ViaCEPResponse, error) {
	url := fmt.Sprintf("%s/ws/%s/json/", c.baseURL, cep)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch location: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return nil, domain.ErrInvalidCEP
	case http.StatusNotFound:
		return nil, domain.ErrCEPNotFound
	case http.StatusOK:
		var location domain.ViaCEPResponse
		if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		if location.Localidade == "" {
			return nil, domain.ErrCEPNotFound
		}

		return &location, nil
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

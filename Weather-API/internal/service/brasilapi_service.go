package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bianavic/fullcycle_desafios/internal/domain"
)

type BrasilAPIService struct{}

func NewBrasilAPIService() *BrasilAPIService {
	return &BrasilAPIService{}
}
func (s *BrasilAPIService) GetLocationByCEP(cep string) (*domain.LocationResponse, error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedLocationData, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, domain.ErrCEPNotFound
	}

	var brasilAPIData struct {
		City     string `json:"city"`
		State    string `json:"state"`
		CEP      string `json:"cep"`
		District string `json:"neighborhood"`
		Street   string `json:"street"`
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &brasilAPIData); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedToParseData, err)
	}

	return &domain.LocationResponse{
		City:     brasilAPIData.City,
		State:    brasilAPIData.State,
		CEP:      brasilAPIData.CEP,
		District: brasilAPIData.District,
		Street:   brasilAPIData.Street,
		Service:  "BrasilAPI",
	}, nil
}

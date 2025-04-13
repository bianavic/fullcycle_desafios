package service

import (
	"log"

	"github.com/bianavic/fullcycle_desafios/internal/domain"
)

type FallbackLocationService struct {
	primary   domain.LocationService
	secondary domain.LocationService
}

func NewFallbackLocationService(primary, secondary domain.LocationService) *FallbackLocationService {
	return &FallbackLocationService{
		primary:   primary,
		secondary: secondary,
	}
}

func (s *FallbackLocationService) GetLocationByCEP(cep string) (*domain.ViaCEPResponse, error) {
	location, err := s.primary.GetLocationByCEP(cep)
	if err == nil && location.Localidade != "" {
		return location, nil
	}

	log.Printf("Primary service failed: %v. Trying fallback...", err)

	location, err = s.secondary.GetLocationByCEP(cep)
	if err != nil {
		log.Printf("Fallback service failed: %v", err)
		return nil, err
	}

	return location, nil
}

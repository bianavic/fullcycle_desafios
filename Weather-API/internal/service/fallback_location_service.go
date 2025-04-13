package service

import (
	"context"
	"time"

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

func (s *FallbackLocationService) GetLocationByCEP(cep string) (*domain.LocationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resultChan := make(chan *domain.LocationResponse, 1)
	errChan := make(chan error, 1)

	go func() {
		loc, err := s.primary.GetLocationByCEP(cep)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- loc
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case <-errChan:
		loc, err := s.secondary.GetLocationByCEP(cep)
		if err != nil {
			return nil, err
		}
		return loc, nil
	case <-ctx.Done():
		loc, err := s.secondary.GetLocationByCEP(cep)
		if err != nil {
			return nil, domain.ErrFailedLocationData
		}
		return loc, nil
	}
}

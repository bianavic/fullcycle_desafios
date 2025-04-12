package repository

import "github.com/bianavic/fullcycle_desafios/internal/domain"

type CEPRepository interface {
	GetLocation(cep string) (*domain.ViaCEPResponse, error)
}

package graph

import (
	"github.com/bianavic/fullcycle_clean-architecture/internal/usecase"
	"github.com/bianavic/fullcycle_clean-architecture/pkg/events"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrderUseCase   usecase.ListOrderUseCase
	EventDispatcher    events.EventDispatcher
}

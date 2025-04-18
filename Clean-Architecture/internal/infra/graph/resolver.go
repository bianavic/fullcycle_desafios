package graph

import (
	"github.com/bianavic/fullcycle_clean-architecture/internal/usecase"
	"github.com/bianavic/fullcycle_clean-architecture/pkg/events"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrderUseCase   usecase.ListOrderUseCase
	EventDispatcher    events.EventDispatcher
}

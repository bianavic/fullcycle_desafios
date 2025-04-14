package service

import (
	"context"
	"github.com/bianavic/fullcycle_clean-architecture/internal/dto"
	"github.com/bianavic/fullcycle_clean-architecture/internal/event"
	"github.com/bianavic/fullcycle_clean-architecture/internal/infra/grpc/pb"
	"github.com/bianavic/fullcycle_clean-architecture/internal/usecase"
	"github.com/bianavic/fullcycle_clean-architecture/pkg/events"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	EventDispatcher    events.EventDispatcher
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, eventDispatcher events.EventDispatcher) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		EventDispatcher:    eventDispatcher,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := dto.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}

	event := event.NewOrderCreated()
	event.SetPayload(dto)
	s.EventDispatcher.Dispatch(event)

	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

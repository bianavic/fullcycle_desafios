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
	ListOrderUseCase   usecase.ListOrderUseCase
	EventDispatcher    events.EventDispatcher
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, listOrderUseCase usecase.ListOrderUseCase, eventDispatcher events.EventDispatcher) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrderUseCase:   listOrderUseCase,
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

func (s *OrderService) ListOrders(_ context.Context, _ *pb.ListOrderRequest) (*pb.ListOrderResponse, error) {
	list := make([]*pb.Order, 0)

	output, err := s.ListOrderUseCase.Execute()
	if err != nil {
		return nil, err
	}

	for _, order := range output {
		dto := &pb.Order{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		}
		list = append(list, dto)
	}

	return &pb.ListOrderResponse{Orders: list}, nil
}

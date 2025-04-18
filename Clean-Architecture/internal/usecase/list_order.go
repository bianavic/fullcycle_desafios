package usecase

import (
	"fmt"
	"github.com/bianavic/fullcycle_clean-architecture/internal/dto"
	"github.com/bianavic/fullcycle_clean-architecture/internal/entity"
	"github.com/bianavic/fullcycle_clean-architecture/pkg/events"
	"log"
)

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewListOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreated events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
		OrderCreated:    OrderCreated,
		EventDispatcher: EventDispatcher,
	}
}

func (c *ListOrderUseCase) Execute() ([]dto.OrderOutputDTO, error) {
	if c.OrderRepository == nil {
		return nil, fmt.Errorf("orderRepository is nil")
	}

	ordersEntity, err := c.OrderRepository.ListOrders()
	if err != nil {
		return nil, err
	}
	log.Printf("orders from repository: %+v", ordersEntity)

	var orders = []dto.OrderOutputDTO{}
	for _, order := range ordersEntity {
		dto := dto.OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.Price + order.Tax,
		}
		orders = append(orders, dto)
	}

	return orders, nil
}

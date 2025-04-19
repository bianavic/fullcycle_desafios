package usecase

import (
	"github.com/bianavic/fullcycle_clean-architecture/internal/dto"
	"github.com/bianavic/fullcycle_clean-architecture/internal/entity"
	"github.com/bianavic/fullcycle_clean-architecture/pkg/events"
	"log"
)

type CreateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface // criar nova ordem no db
	OrderCreated    events.EventInterface           // o evento que será disparado
	EventDispatcher events.EventDispatcherInterface // dispara o evento
}

func NewCreateOrderUseCase(
	OrderRepositoryInterface entity.OrderRepositoryInterface,
	OrderCreated events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: OrderRepositoryInterface,
		OrderCreated:    OrderCreated,
		EventDispatcher: EventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input dto.OrderInputDTO) (dto.OrderOutputDTO, error) {
	log.Printf("CreateOrder input - ID: %s, Price: %f, Tax: %f", input.ID, input.Price, input.Tax)

	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}

	err := order.CalculateFinalPrice()
	if err != nil {
		return dto.OrderOutputDTO{}, err
	}
	if err := c.OrderRepository.Save(&order); err != nil {
		log.Printf("error saving order: %v", err)
		return dto.OrderOutputDTO{}, err
	}

	dto := dto.OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}

	c.OrderCreated.SetPayload(dto)
	c.EventDispatcher.Dispatch(c.OrderCreated)

	return dto, nil
}

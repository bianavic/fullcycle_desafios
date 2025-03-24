package usecase

import (
	"database/sql"
	"time"

	"github.com/bianavic/fullcycle_clean-architecture/internal/dto"
	"github.com/bianavic/fullcycle_clean-architecture/internal/entity"
)

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrderUseCase(OrderRepository entity.OrderRepositoryInterface) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *ListOrderUseCase) Execute() ([]dto.OrderOutputDTO, error) {
	var orders = []dto.OrderOutputDTO{}

	ordersEntity, err := c.OrderRepository.List()
	if err != nil {
		return nil, err
	}

	location, _ := time.LoadLocation("Local")

	for _, order := range ordersEntity {
		dto := dto.OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.Price + order.Tax,
			CreatedAt:  convertToTimezone(order.CreatedAt, location),
		}
		orders = append(orders, dto)
	}

	return orders, nil
}

func convertToTimezone(t sql.NullTime, location *time.Location) string {
	if t.Valid {
		return t.Time.In(location).Format("2006-01-02 15:04:05 -07:00")
	}
	return ""
}

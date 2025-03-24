package usecase

import (
	"github.com/bianavic/fullcycle_clean-architecture/internal/entity"
	"github.com/bianavic/fullcycle_clean-architecture/pkg/events"
)

// USE CASE:
// 1. RECEBE os dados - DTO - é um input do cliente - anemica
// 2. CALCULA final price
// 3. SALVA no banco
// 4. cria um evento
// 5. DISPARA o evento no event dispatcher

type OrderInputDTO struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

// retorna DTO - dados para o cliente com o final price
type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

// USE CASE com 3 compomentes, mapeados como INTERFACES - Inversão de Controle
type CreateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface // criar nova ordem no db
	OrderCreated    events.EventInterface           // o evento que será disparado
	EventDispatcher events.EventDispatcherInterface // dispara o evento
}

// Cria o Use Case
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

// Executa o Use Case: recebe DTO, calcula o final price, salva no banco e dispara o evento
// é o metodo PRINCIPAL do use case que recebe como parametro o DTO,
// transforma o DTO na ordem
// executa a regra de negocio (o calculo do final price)
// e retorna o DTO
func (c *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}
	order.CalculateFinalPrice()

	// salva no banco de dados e o use case nao sabe como é feito
	if err := c.OrderRepository.Save(&order); err != nil {
		return OrderOutputDTO{}, err
	}

	// preparar o output DTO
	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}

	// passa o DTO do output e passa para o evento para a interface (SetPayload) - o evento tem que ter um payload
	c.OrderCreated.SetPayload(dto)
	// dispara o evento para uma interface
	c.EventDispatcher.Dispatch(c.OrderCreated)

	return dto, nil
}

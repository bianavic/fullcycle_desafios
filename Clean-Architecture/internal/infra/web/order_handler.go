package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bianavic/fullcycle_clean-architecture/internal/dto"
	"github.com/bianavic/fullcycle_clean-architecture/internal/entity"
	usecase "github.com/bianavic/fullcycle_clean-architecture/internal/usecase"
	"github.com/bianavic/fullcycle_clean-architecture/pkg/events"
)

type WebOrderHandler struct {
	EventDispatcher   events.EventDispatcherInterface
	OrderRepository   entity.OrderRepositoryInterface
	OrderCreatedEvent events.EventInterface
}

func NewWebOrderHandler(
	EventDispatcher events.EventDispatcherInterface,
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreatedEvent events.EventInterface,
) *WebOrderHandler {
	return &WebOrderHandler{
		EventDispatcher:   EventDispatcher,
		OrderRepository:   OrderRepository,
		OrderCreatedEvent: OrderCreatedEvent,
	}
}

func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	// validation for nil dependencies
	if h == nil || h.OrderRepository == nil || h.EventDispatcher == nil || h.OrderCreatedEvent == nil {
		http.Error(w, "Handler not properly initialized", http.StatusInternalServerError)
		return
	}

	log.Println("Received order creation request")

	var dto dto.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebOrderHandler) List(w http.ResponseWriter, r *http.Request) {
	if h.OrderRepository == nil {
		http.Error(w, "OrderRepository is nil", http.StatusInternalServerError)
		return
	}
	
	listOrder := usecase.NewListOrderUseCase(h.OrderRepository)
	output, err := listOrder.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var orders []dto.OrderOutputDTO
	for _, order := range output {
		orders = append(orders, dto.OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
			CreatedAt:  order.CreatedAt,
		})
	}

	if len(orders) == 0 {
		orders = []dto.OrderOutputDTO{}
	}

	response := dto.OrdersOutputDTO{
		Orders: orders,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

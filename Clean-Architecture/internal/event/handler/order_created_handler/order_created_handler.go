package order_created_handler

import (
	"encoding/json"
	"fmt"
	"github.com/bianavic/fullcycle_clean-architecture/pkg/events"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

type OrderCreatedHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewOrderCreatedHandler(rabbitMQChannel *amqp.Channel) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

// é o que permite impleemntar a interface para passar no event dispatcher
func (h *OrderCreatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Order created: %v\n", event.GetPayload())

	jsonOutput, err := json.Marshal(event.GetPayload()) // pega o payload do evento
	if err != nil {
		log.Printf("Error marshaling event payload: %v", err)
		return
	}

	msgRabbitmq := amqp.Publishing{ //  e transforma em json para enviar para o rabbitmq
		ContentType: "application/json",
		Body:        jsonOutput, // payload do evento
	}

	err = h.RabbitMQChannel.Publish(
		"amq.direct", // exchange
		"",           // key name
		false,        // mandatory
		false,        // immediate
		msgRabbitmq,  // message to publish (payload do evento / json)
	)
	if err != nil {
		log.Printf("Error publishing message: %v", err)
	}
}

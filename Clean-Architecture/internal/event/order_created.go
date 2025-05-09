package event

// PODE ESTAR NA PASTA INFREA OU NA PASTA ENTITY

// EventInterface é a interface que define um evento
// qdo dou dispatch em um evento, ele executa varios handler
type OrderCreated struct {
	Name    string
	Payload interface{}
}

func NewOrderCreated() *OrderCreated {
	return &OrderCreated{
		Name: "OrderCreated",
	}
}

func (e *OrderCreated) GetName() string {
	return e.Name
}

func (e *OrderCreated) GetPayload() interface{} {
	return e.Payload
}

func (e *OrderCreated) SetPayload(payload interface{}) {
	e.Payload = payload
}

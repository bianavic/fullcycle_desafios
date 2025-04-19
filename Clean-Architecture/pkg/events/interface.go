package events

import (
	"sync"
)

type EventInterface interface {
	GetName() string
	GetPayload() interface{}
	SetPayload(payload interface{})
}

type EventHandlerInterface interface {
	Handle(event EventInterface, wg *sync.WaitGroup)
}

type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface) error
	Dispatch(event EventInterface) error
	Remove(eventName string, handler EventHandlerInterface) error
	Has(eventName string, handler EventHandlerInterface) bool
	Clear()
}

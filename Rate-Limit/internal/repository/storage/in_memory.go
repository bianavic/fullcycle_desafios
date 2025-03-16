package storage

import (
	"context"
	"sync"
	"time"
)

// InMemory implements the StorageStrategy interface using an in-memory map.
type InMemory struct {
	mu    sync.Mutex
	store map[string]int
}

// NewInMemory creates a new InMemory instance.
func NewInMemory() *InMemory {
	return &InMemory{
		store: make(map[string]int),
	}
}

// Increment increments the value for a key and returns the new value.
func (i *InMemory) Increment(ctx context.Context, key string, expiration time.Duration) (int, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.store[key]++
	return i.store[key], nil
}

// Get retrieves the value for a key.
func (i *InMemory) Get(ctx context.Context, key string) (int, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	val, exists := i.store[key]
	if !exists {
		return 0, nil
	}
	return val, nil
}

// Set sets the value for a key with an expiration time.
func (i *InMemory) Set(ctx context.Context, key string, value int, expiration time.Duration) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.store[key] = value
	return nil
}

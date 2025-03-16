package limiter

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInMemory_IncrementIncreasesValue(t *testing.T) {
	store := NewInMemory()
	ctx := context.Background()

	val, err := store.Increment(ctx, "key1", time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	val, err = store.Increment(ctx, "key1", time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, 2, val)
}

func TestInMemory_GetReturnsValue(t *testing.T) {
	store := NewInMemory()
	ctx := context.Background()

	_, err := store.Increment(ctx, "key1", time.Minute)
	assert.NoError(t, err)

	val, err := store.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
}

func TestInMemory_GetReturnsZeroForNonExistentKey(t *testing.T) {
	store := NewInMemory()
	ctx := context.Background()

	val, err := store.Get(ctx, "nonexistent")
	assert.NoError(t, err)
	assert.Equal(t, 0, val)
}

func TestInMemory_SetStoresValue(t *testing.T) {
	store := NewInMemory()
	ctx := context.Background()

	err := store.Set(ctx, "key1", 5, time.Minute)
	assert.NoError(t, err)

	val, err := store.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.Equal(t, 5, val)
}

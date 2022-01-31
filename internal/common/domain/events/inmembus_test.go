package events

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInMemBus_Subscribe(t *testing.T) {
	bus := NewInMemoryEventBus()

	handlerA := NewFuncHandler(func(ctx context.Context, event Event) error {
		return nil
	}, func(event Event) bool {
		return true
	})

	handlerB := NewFuncHandler(func(ctx context.Context, event Event) error {
		return nil
	}, func(event Event) bool {
		return true
	})

	bus.Subscribe(handlerA, handlerA, handlerB)
	require.Equal(t, 3, len(bus.handlers))

	bus.Unsubscribe(handlerA)
	require.Equal(t, 1, len(bus.handlers))

	bus.Unsubscribe(handlerB)
	require.Equal(t, 0, len(bus.handlers))
}

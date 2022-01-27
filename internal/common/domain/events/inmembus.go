package events

import (
	"context"
)

// EventBus is an in-memory implementation of the events.Bus.
type EventBus struct {
	handlers []Handler
}

// NewInMemoryEventBus initializes a new EventBus.
func NewInMemoryEventBus() *EventBus {
	return &EventBus{
		handlers: make([]Handler, 0),
	}
}

// Publish implements the events.Bus interface.
func (b *EventBus) Publish(ctx context.Context, events []Event) error {
	for _, evt := range events {
		for _, handler := range b.handlers {
			if true == handler.IsSubscribeTo(evt) {
				_ = handler.Handle(ctx, evt)
			}
		}
	}

	return nil
}

// Subscribe implements the events.Bus interface.
func (b *EventBus) Subscribe(handler Handler) {
	b.handlers = append(b.handlers, handler)
}

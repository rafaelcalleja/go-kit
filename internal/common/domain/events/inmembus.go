package events

import (
	"context"
)

// EventBus is an in-memory implementation of the events.Bus.
type EventBus struct {
	handlers map[Type][]Handler
}

// NewInMemoryEventBus initializes a new EventBus.
func NewInMemoryEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[Type][]Handler),
	}
}

// Publish implements the events.Bus interface.
func (b *EventBus) Publish(ctx context.Context, events []Event) error {
	for _, evt := range events {
		handlers, ok := b.handlers[evt.Type()]
		if !ok {
			return nil
		}

		for _, handler := range handlers {
			_ = handler.Handle(ctx, evt)
		}
	}

	return nil
}

// Subscribe implements the events.Bus interface.
func (b *EventBus) Subscribe(evtType Type, handler Handler) {
	subscribersForType, ok := b.handlers[evtType]
	if !ok {
		b.handlers[evtType] = []Handler{handler}
	}

	b.handlers[evtType] = append(subscribersForType, handler)
}

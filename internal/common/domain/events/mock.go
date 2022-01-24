package events

import (
	"context"
)

// EventBusMock is a mock implementation of the events.Bus.
type EventBusMock struct {
	PublishFn   func(ctx context.Context, events []Event) error
	SubscribeFn func(evtType Type, handler Handler)
}

// NewMockEventBus initializes a new EventBus.
func NewMockEventBus() *EventBusMock {
	return &EventBusMock{
		PublishFn:   func(ctx context.Context, events []Event) error { return nil },
		SubscribeFn: func(evtType Type, handler Handler) {},
	}
}

// Publish implements the events.Bus interface.
func (b *EventBusMock) Publish(ctx context.Context, events []Event) error {
	return b.PublishFn(ctx, events)
}

// Subscribe implements the events.Bus interface.
func (b *EventBusMock) Subscribe(evtType Type, handler Handler) {
	b.SubscribeFn(evtType, handler)
}

const MockEventType Type = "data.mocks"

type MockEvent struct {
	BaseEvent
}

func (e MockEvent) Type() Type {
	return MockEventType
}

func NewMockEvent(id string) MockEvent {
	return MockEvent{
		BaseEvent: NewBaseEvent(id),
	}
}

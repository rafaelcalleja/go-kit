package domain

import (
	"github.com/rafaelcalleja/go-kit/internal/common/domain/events"
)

const ProductCreatedEventType events.Type = "events.product.created"

type ProductCreatedEvent struct {
	events.BaseEvent
	id string
}

func NewProductCreatedEvent(id string) ProductCreatedEvent {
	return ProductCreatedEvent{
		id: id,

		BaseEvent: events.NewBaseEvent(id),
	}
}

func (e ProductCreatedEvent) Type() events.Type {
	return ProductCreatedEventType
}

func (e ProductCreatedEvent) ProductId() string {
	return e.id
}

package events

import (
	"context"
)

type Bus interface {
	// Publish is the method used to publish new data.
	Publish(context.Context, []Event) error
	// Subscribe is the method used to subscribe new event handlers.
	Subscribe(...*Handler)
	Unsubscribe(...*Handler)
}

// Handler defines the expected behaviour from an event handler.
type Handler interface {
	Handle(context.Context, Event) error
	IsSubscribeTo(Event) bool
}

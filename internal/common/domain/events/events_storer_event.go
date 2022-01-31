package events

import (
	"context"
)

type StoreEventsOnEventCreated struct {
	store Store
}

func NewStoreEventsOnEventCreated(store Store) *Handler {
	var handler Handler = StoreEventsOnEventCreated{
		store: store,
	}

	return &handler
}

func (e StoreEventsOnEventCreated) Handle(_ context.Context, evt Event) error {
	e.store.Append(evt)

	return nil
}

func (e StoreEventsOnEventCreated) IsSubscribeTo(Event) bool {
	return true
}

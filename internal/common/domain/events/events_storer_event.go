package events

import (
	"context"
)

type StoreEventsOnEventCreated struct {
	store Store
}

func NewStoreEventsOnEventCreated(store Store) StoreEventsOnEventCreated {
	return StoreEventsOnEventCreated{
		store: store,
	}
}

func (e StoreEventsOnEventCreated) Handle(_ context.Context, evt Event) error {
	e.store.Append(evt)

	return nil
}

func (e StoreEventsOnEventCreated) IsSubscribeTo(Event) bool {
	return true
}

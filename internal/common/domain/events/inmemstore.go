package events

import (
	"sync"
)

type EventStoreInMem struct {
	mu      sync.Mutex
	events  []Event
	indexes map[string]int
}

func NewInMemEventStore() *EventStoreInMem {
	return NewInMemEventStoreWith()
}

func NewInMemEventStoreWith(options ...func(*EventStoreInMem)) *EventStoreInMem {
	var eventStore = new(EventStoreInMem)

	eventStore.mu = sync.Mutex{}
	eventStore.events = make([]Event, 0)
	eventStore.indexes = make(map[string]int)

	for _, option := range options {
		option(eventStore)
	}

	return eventStore
}

func WithMemIndex(indexes map[string]int) func(*EventStoreInMem) {
	return func(m *EventStoreInMem) {
		m.indexes = indexes
	}
}

func WithMemStore(events []Event) func(*EventStoreInMem) {
	return func(m *EventStoreInMem) {
		m.events = events
	}
}

func (e *EventStoreInMem) Events() []Event {
	return e.events
}

func (e *EventStoreInMem) Append(event Event) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.events = append(e.events, event)
	e.indexes[event.ID()] = len(e.events) - 1
}

func (e *EventStoreInMem) AllStoredEventsSince(eventId EventId) []Event {
	e.mu.Lock()
	defer e.mu.Unlock()

	if val, ok := e.indexes[eventId.String()]; ok {
		return e.events[val:]
	}

	return make([]Event, 0)
}

func (e *EventStoreInMem) StoredEventsSince(eventId EventId, limit int) []Event {
	return e.AllStoredEventsSince(eventId)[:limit]
}

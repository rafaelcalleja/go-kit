package events

import (
	"sync"
)

type EventStoreOriginator struct {
	mu         sync.Mutex
	originator originator
	indexes    map[string]int
}

func NewEventStoreOriginator() EventStoreOriginator {
	return NewEventStoreOriginatorWith()
}

func NewEventStoreOriginatorWith(options ...func(*EventStoreOriginator)) EventStoreOriginator {
	var eventStore = new(EventStoreOriginator)

	eventStore.mu = sync.Mutex{}
	eventStore.originator = originator{
		events: make([]Event, 0),
	}

	eventStore.indexes = make(map[string]int)

	for _, option := range options {
		option(eventStore)
	}

	return *eventStore
}

func (e *EventStoreOriginator) Events() []Event {
	return e.originator.events
}

func (e *EventStoreOriginator) Append(event Event) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.originator.events = append(e.originator.events, event)
	e.indexes[event.ID()] = len(e.originator.events) - 1
}

func (e *EventStoreOriginator) AllStoredEventsSince(eventId EventId) []Event {
	e.mu.Lock()
	defer e.mu.Unlock()

	if val, ok := e.indexes[eventId.String()]; ok {
		return e.originator.events[val:]
	}

	return make([]Event, 0)
}

func (e *EventStoreOriginator) StoredEventsSince(eventId EventId, limit int) []Event {
	return e.AllStoredEventsSince(eventId)[:limit]
}

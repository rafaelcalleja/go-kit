package events

import (
	"sync"
)

type EventStoreInMem struct {
	mu      sync.Mutex
	events  []Event
	indexes map[string]int
}

func NewEventStoreInMem() EventStoreInMem {
	return EventStoreInMem{
		events:  make([]Event, 0),
		indexes: make(map[string]int),
	}
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

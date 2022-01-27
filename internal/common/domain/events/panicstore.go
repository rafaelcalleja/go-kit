package events

type PanicStore struct {
	events  []Event
	indexes map[string]int
}

func NewPanicStore() PanicStore {
	return NewPanicStoreWith()
}

func NewPanicStoreWith(options ...func(*PanicStore)) PanicStore {
	var eventStore = new(PanicStore)

	eventStore.events = make([]Event, 0)
	eventStore.indexes = make(map[string]int)

	for _, option := range options {
		option(eventStore)
	}

	return *eventStore
}

func (e *PanicStore) Events() []Event {
	panic("err")
}

func (e *PanicStore) Append(event Event) {
}

func (e *PanicStore) AllStoredEventsSince(eventId EventId) []Event {
	return make([]Event, 0)
}

func (e *PanicStore) StoredEventsSince(eventId EventId, limit int) []Event {
	return e.AllStoredEventsSince(eventId)[:limit]
}

package events

type originator struct {
	events []Event
}

func (o *originator) createMemento(events []Event) *memento {
	clone := make([]interface{}, len(events))
	for i := range events {
		clone[i] = events[i]
	}
	return &memento{
		data: clone,
	}
}

func (o *originator) restoreMemento(m *memento) {
	currEvents := m.getData()
	o.events = make([]Event, len(currEvents))
	for i := range currEvents {
		o.events[i] = currEvents[i].(Event)
	}
}

func (o *originator) setEvents(events []Event) {
	o.events = events
}

func (o *originator) getEvents() []Event {
	return o.events
}

package events

type originator struct {
	events []Event
}

func (o *originator) createMemento() *memento {
	clone := make([]interface{}, len(o.events))
	for i := range o.events {
		clone[i] = o.events[i]
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

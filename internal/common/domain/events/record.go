package events

type RecordEvents interface {
	Record(evt Event)
	PullEvents() []Event
}

type ImplementRecordableEvents struct {
	events []Event
}

func (i *ImplementRecordableEvents) Record(event Event) {
	i.events = append(i.events, event)
}

func (i *ImplementRecordableEvents) PullEvents() []Event {
	events := i.events
	i.events = []Event{}

	return events
}

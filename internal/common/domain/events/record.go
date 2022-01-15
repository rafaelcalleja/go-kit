package events

type RecordEvents interface {
	Record(evt Event)
	PullEvents() []Event
}

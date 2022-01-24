package events

type Store interface {
	Append(event Event)
	AllStoredEventsSince(eventId EventId) []Event
	StoredEventsSince(eventId EventId, limit int) []Event
}

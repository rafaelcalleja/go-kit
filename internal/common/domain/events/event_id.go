package events

import (
	"fmt"

	"github.com/rafaelcalleja/go-kit/uuid"
)

type EventId struct {
	id uuid.UUID
}

type AggregateId struct {
	EventId
}

func NewAggregateId(id string) (AggregateId, error) {
	idVO, err := NewEventId(id)

	if nil != err {
		return AggregateId{}, err
	}

	return AggregateId{idVO}, nil
}

func NewEventId(id string) (EventId, error) {
	idVO, err := uuid.New().Parse(id)

	if nil != err {
		return EventId{}, fmt.Errorf("%w: %s", ErrWrongUuid, id)
	}

	return EventId{idVO}, nil
}

func (e EventId) String() string {
	return uuid.New().String(e.id)
}

func (e EventId) Equals(other EventId) bool {
	return other.String() == e.String()
}

package events

import (
	"errors"
	"time"

	"github.com/rafaelcalleja/go-kit/uuid"
)

var (
	ErrWrongUuid = errors.New("wrong uuid")
)

// Type represents a domain data type.
type Type string

// Event represents a domain command.
type Event interface {
	ID() string
	AggregateID() string
	OccurredOn() time.Time
	Type() Type
}

type BaseEvent struct {
	eventID     EventId
	aggregateID AggregateId
	occurredOn  time.Time
}

func NewBaseEvent(aggregateID string) BaseEvent {
	eventIdVo, _ := NewEventId(uuid.New().String(uuid.New().Create()))
	aggregateIDVO, _ := NewAggregateId(aggregateID)

	return BaseEvent{
		eventID:     eventIdVo,
		aggregateID: aggregateIDVO,
		occurredOn:  time.Now(),
	}
}

func (b BaseEvent) ID() string {
	return b.eventID.String()
}

func (b BaseEvent) OccurredOn() time.Time {
	return b.occurredOn
}

func (b BaseEvent) AggregateID() string {
	return b.aggregateID.String()
}

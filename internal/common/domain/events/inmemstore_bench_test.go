package events

import (
	"fmt"
	"github.com/rafaelcalleja/go-kit/uuid"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

var (
	eventStore = NewInMemEventStore()
	events     = make([]Event, 0)
)

func BenchmarkNewInMemoryEventBus(b *testing.B) {
	rand.Seed(time.Now().Unix())

	for a := 0; a <= 4; a++ {
		appendEvents(1000000)
		b.Run(fmt.Sprintf("input_size_%d", len(events)), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				findEventId := rand.Intn(len(events))
				eventToFind := events[findEventId]

				eventToFindVO, _ := NewEventId(eventToFind.ID())
				newEventsSlice := eventStore.AllStoredEventsSince(eventToFindVO)

				eventToFondVO, _ := NewEventId(newEventsSlice[0].ID())

				assert.True(b, eventToFindVO.Equals(eventToFondVO))
				assert.GreaterOrEqual(b, len(newEventsSlice), 1)
			}
		})
	}
}

func appendEvents(numEvents int) {
	for a := 0; a < numEvents; a++ {
		event := NewMockEvent(uuid.New().String(uuid.New().Create()))
		events = append(events, event)
		eventStore.Append(event)
	}
}

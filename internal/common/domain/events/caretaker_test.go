package events

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCaretaker(t *testing.T) {
	events := make([]Event, 0)
	caretaker := make(CareTaker, 0)
	originator := &originator{
		events,
	}

	require.Equal(t, events, originator.getEvents())
	require.NotSame(t, events, originator.getEvents())
	caretaker.Add(originator.createMemento())

	events = append(events, NewMockEvent("1"))
	eventsA := events
	originator.setEvents(events)
	require.Equal(t, events, originator.getEvents())
	require.NotSame(t, events, originator.getEvents())

	events = append(events, NewMockEvent("2"))
	eventsB := events
	caretaker.Add(originator.createMemento())
	originator.setEvents(events)
	require.Equal(t, events, originator.getEvents())
	require.NotSame(t, events, originator.getEvents())

	caretaker.Add(originator.createMemento())

	originator.restoreMemento(caretaker.Pop())
	restore1 := originator.getEvents()
	require.Equal(t, eventsB, restore1)
	require.NotSame(t, events, restore1)

	originator.restoreMemento(caretaker.Pop())
	restore2 := originator.getEvents()
	require.Equal(t, eventsA, restore2)
	require.NotEqual(t, restore1, restore2)

	originator.restoreMemento(caretaker.Pop())
	restore3 := originator.getEvents()
	require.Equal(t, make([]Event, 0), restore3)
}

func TestIteration(t *testing.T) {
	events := make([]Event, 0)
	originator := &originator{
		events,
	}

	mementoA, mementoB := originator.createMemento(), originator.createMemento()

	ct := make(CareTaker, 0)
	ct.Add(mementoA, mementoB)
	require.Equal(t, 2, len(ct))

	m1 := ct.Pop()
	require.Equal(t, 1, len(ct))
	require.Same(t, mementoB, m1)

	m2 := ct.Pop()
	require.Equal(t, 0, len(ct))
	require.Same(t, mementoA, m2)

	require.NotSame(t, m2, m1)
}

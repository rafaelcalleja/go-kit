package events

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestInMemTransaction_Rollback(t *testing.T) {
	session := NewInMemEventTransaction(nil)
	mockEvent := NewMockEvent("A")
	expected := session.GetEvents()

	tx, _ := session.Begin()

	session.AddEvent(mockEvent)

	_ = tx.Rollback()

	require.Equal(t, expected, session.GetEvents())
	require.Equal(t, 0, len(session.GetEvents()))
}

func TestInMemTransaction_Success(t *testing.T) {
	session := NewInMemEventTransaction(nil)
	mockEvent := NewMockEvent("A")

	tx, _ := session.Begin()

	session.AddEvent(mockEvent)
	expected := session.GetEvents()

	_ = tx.Commit()

	require.Equal(t, expected, session.GetEvents())
	require.Equal(t, 1, len(session.GetEvents()))
}

func TestInMemTransaction_MultipleTransaction_Rollback(t *testing.T) {
	session := NewInMemEventTransaction(nil)
	mockEventExpected := NewMockEvent("A")
	mockEventCancel := NewMockEvent("B")

	tx, _ := session.Begin()

	session.AddEvent(mockEventExpected)
	expected := session.GetEvents()

	tx2, _ := session.Begin()
	session.AddEvent(mockEventCancel)

	_ = tx2.Rollback()
	_ = tx.Commit()

	require.Equal(t, expected, session.GetEvents())
	require.Equal(t, 1, len(session.GetEvents()))
	require.Equal(t, mockEventExpected.ID(), session.GetEvents()[0].ID())
}

func TestInMemTransaction_MultipleTransaction_Success(t *testing.T) {
	session := NewInMemEventTransaction(nil)
	mockEvent := NewMockEvent("A")

	tx, _ := session.Begin()
	session.AddEvent(mockEvent)

	tx2, _ := session.Begin()
	session.AddEvent(mockEvent)

	expected := session.GetEvents()

	_ = tx2.Commit()
	_ = tx.Commit()

	require.Equal(t, expected, session.GetEvents())
	require.Equal(t, 2, len(session.GetEvents()))
}

func TestInMemTransaction_Concurrent_Success(t *testing.T) {
	events := []Event{NewMockEvent("00000000-0000-0000-0000-000000000001"), NewMockEvent("00000000-0000-0000-0000-000000000002"), NewMockEvent("00000000-0000-0000-0000-000000000003")}
	WrongEvent := NewMockEvent("00000000-0000-0000-0000-000000000099")

	expected := []string{
		"00000000-0000-0000-0000-000000000001",
		"00000000-0000-0000-0000-000000000002",
		"00000000-0000-0000-0000-000000000003",
		"00000000-0000-0000-0000-000000000099",
	}

	session := NewInMemEventTransaction(nil)
	tx, _ := session.Begin()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		tx2, _ := session.Begin()
		session.AddEvent(WrongEvent)
		_ = tx2.Commit()
		defer wg.Done()
	}()

	go func() {
		time.Sleep(1 * time.Second)
		session.AddEvent(events...)
		_ = tx.Commit()
		defer wg.Done()
	}()
	wg.Wait()

	actual := make([]string, 0)

	for _, v := range session.GetEvents() {
		actual = append(actual, v.AggregateID())
	}

	require.Equal(t, expected, actual)
}

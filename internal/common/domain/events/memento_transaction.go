package events

import (
	"fmt"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
	"sync"
)

type MementoTx struct {
	originator *originator
	memento    *memento
	store      *EventStoreInMem
	mu         sync.Mutex
}

func NewMementoTx(store *EventStoreInMem) *MementoTx {
	events := store.events

	for k, v := range store.events {
		events[k] = v
	}

	originator := &originator{
		events,
	}

	return &MementoTx{
		originator: originator,
		store:      store,
	}
}

func (o *MementoTx) Begin() (transaction.Transaction, error) {
	o.mu.Lock()
	o.memento = o.originator.createMemento()

	return transaction.Transaction(o), nil
}

func (o *MementoTx) Rollback() error {
	defer o.mu.Unlock()
	o.originator.restoreMemento(o.memento)
	o.memento = nil

	o.store.events = make([]Event, 0)

	fmt.Println("ROOOOOOOOOOOOOOOOOL", len(o.originator.events), len(o.store.events), o.store, len(o.store.Events()))

	for _, v := range o.originator.events {
		o.store.Append(v)
	}

	return nil
}

func (o *MementoTx) Commit() error {
	defer o.mu.Unlock()
	return nil
}

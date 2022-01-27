package events

import (
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type MementoTx struct {
	originator *originator
	memento    *memento
	store      *EventStoreInMem
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
	o.memento = o.originator.createMemento()

	return transaction.Transaction(o), nil
}

func (o *MementoTx) Rollback() error {
	o.originator.restoreMemento(o.memento)

	o.store.events = make([]Event, 0)

	for _, v := range o.originator.events {
		o.store.Append(v)
	}

	return nil
}

func (o *MementoTx) Commit() error {
	return nil
}

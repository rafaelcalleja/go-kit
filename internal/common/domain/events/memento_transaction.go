package events

import (
	"context"
	"sync"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type MementoTx struct {
	backup []Event
	store  *EventStoreInMem
	mu     sync.Mutex
	mux    sync.Mutex
}

func NewMementoTx(store *EventStoreInMem) *MementoTx {
	m := &MementoTx{
		store: store,
	}

	m.doBackup()

	return m
}

func (o *MementoTx) doBackup() {
	o.mux.Lock()
	defer o.mux.Unlock()

	backup := make([]Event, len(o.store.events))

	copy(backup, o.store.events)

	o.backup = backup
}

func (o *MementoTx) doRestore() {
	o.mux.Lock()
	defer o.mux.Unlock()

	restore := make([]Event, len(o.backup))

	copy(restore, o.backup)

	o.backup = make([]Event, 0)

	o.store.events = make([]Event, 0)

	for _, v := range restore {
		o.store.Append(v)
	}

	o.store.events = restore
}

func (o *MementoTx) Begin(_ context.Context) (transaction.Transaction, error) {
	o.mu.Lock()
	o.doBackup()

	return transaction.Transaction(o), nil
}

func (o *MementoTx) Rollback() error {
	defer o.mu.Unlock()
	o.doRestore()
	return nil
}

func (o *MementoTx) Commit() error {
	defer o.mu.Unlock()
	return nil
}

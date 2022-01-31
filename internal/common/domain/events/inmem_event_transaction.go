package events

import (
	"bytes"
	"runtime/debug"
	"strconv"
	"sync"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type InMemEventTransaction struct {
	mu         sync.Mutex
	caretaker  CareTaker
	originator *originator
	memento    *memento
	locker     map[int]int
	curGid     int
}

func NewInMemEventTransaction(events []Event) InMemEventTransaction {
	if nil == events {
		events = make([]Event, 0)
	}

	caretaker := make(CareTaker, 0)
	originator := &originator{
		events,
	}

	return InMemEventTransaction{
		caretaker:  caretaker,
		originator: originator,
		locker:     make(map[int]int),
		curGid:     -1,
	}
}

func (o *InMemEventTransaction) AddEvent(event ...Event) {
	o.originator.events = append(o.originator.events, event...)
}

func (o *InMemEventTransaction) GetEvents() []Event {
	return o.originator.getEvents()
}

func (o *InMemEventTransaction) gid() int {
	gr := bytes.Fields(debug.Stack())[1]
	gid, _ := strconv.Atoi(string(gr))

	return gid
}

func (o *InMemEventTransaction) Begin() (transaction.Transaction, error) {
	gid := o.gid()

	if _, ok := o.locker[gid]; !ok {
		o.mu.Lock()
		o.curGid = gid
	}

	o.locker[gid] = o.locker[gid] + 1

	o.memento = o.originator.createMemento(nil)

	return transaction.Transaction(o), nil
}

func (o *InMemEventTransaction) Rollback() error {
	defer o.endTransaction()
	o.originator.restoreMemento(o.memento)
	o.memento = nil
	return nil
}

func (o *InMemEventTransaction) Commit() error {
	defer o.endTransaction()

	return nil
}

func (o *InMemEventTransaction) endTransaction() {
	gid := o.curGid
	if _, ok := o.locker[gid]; !ok {
		panic("no puedes Commit/Rollback")
	}

	if val, ok := o.locker[gid]; ok && val > 1 {
		o.locker[gid] = o.locker[gid] - 1
		return
	}

	delete(o.locker, gid)
	o.curGid = -1
	o.mu.Unlock()
}

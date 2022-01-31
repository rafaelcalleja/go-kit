package adapters

import (
	"database/sql"
	"sync"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type TransactionInitializerExecutorDb struct {
	db       *sql.DB
	executor transaction.Executor
	stack    []interface{}
	mu       sync.Mutex
}

func NewTransactionInitializerExecutorDb(db *sql.DB, executor transaction.Executor) *TransactionInitializerExecutorDb {
	if db == nil {
		panic("missing db")
	}

	executor.WithConnection(db)

	stack := emptyStack(db)

	return &TransactionInitializerExecutorDb{
		db:       db,
		executor: executor,
		stack:    stack,
	}
}

func emptyStack(db *sql.DB) []interface{} {
	stack := make([]interface{}, 0)
	stack = append(stack, db)

	return stack
}

func (i *TransactionInitializerExecutorDb) last() interface{} {
	return i.stack[len(i.stack)-1]
}

func (i *TransactionInitializerExecutorDb) pop() interface{} {
	if len(i.stack) == 1 {
		return i.stack[0]
	}

	n := len(i.stack) - 1

	elem := i.stack[n]

	i.stack = i.stack[:n]

	return elem
}

func (i *TransactionInitializerExecutorDb) Begin() (transaction.Transaction, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	tx, err := i.db.Begin()

	i.stack = append(i.stack, tx)
	i.executor.WithConnection(i.last().(transaction.Connection))

	return i, err
}

func (i *TransactionInitializerExecutorDb) Rollback() error {
	i.mu.Lock()
	defer func() {
		defer i.mu.Unlock()
		err := recover()
		if nil != err {
			i.stack = emptyStack(i.db)
			return
		}
		i.pop() //Discard current tx
		i.executor.WithConnection(i.pop().(transaction.Connection))
	}()

	return i.last().(transaction.Transaction).Rollback()
}

func (i *TransactionInitializerExecutorDb) Commit() error {
	i.mu.Lock()
	defer func() {
		defer i.mu.Unlock()
		err := recover()
		if nil != err {
			i.stack = emptyStack(i.db)
			return
		}
		i.pop() //Discard current tx
		i.executor.WithConnection(i.pop().(transaction.Connection))
	}()

	return i.last().(transaction.Transaction).Commit()
}

package adapters

import (
	"database/sql"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
	"sync"
)

type TransactionInitializerExecutorSimpleDb struct {
	db       *sql.DB
	executor transaction.Executor
	mu       sync.Mutex
}

func NewTransactionInitializerExecutorSimpleDb(db *sql.DB, executor transaction.Executor) *TransactionInitializerExecutorSimpleDb {
	if db == nil {
		panic("missing db")
	}

	return &TransactionInitializerExecutorSimpleDb{
		db:       db,
		executor: executor,
	}
}

func (i *TransactionInitializerExecutorSimpleDb) Begin() (transaction.Transaction, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if err := i.db.Ping(); err != nil {
		return nil, err
	}

	tx, err := i.db.Begin()
	if err != nil {
		return nil, err
	}

	i.executor.WithConnection(tx)

	return i, err
}

func (i *TransactionInitializerExecutorSimpleDb) Rollback() (err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	defer i.executor.WithConnection(i.db)

	return i.executor.GetConnection().(transaction.Transaction).Rollback()
}

func (i *TransactionInitializerExecutorSimpleDb) Commit() error {
	i.mu.Lock()
	defer i.mu.Unlock()
	defer i.executor.WithConnection(i.db)

	return i.executor.GetConnection().(transaction.Transaction).Commit()
}

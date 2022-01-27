package adapters

import (
	"database/sql"
	"github.com/rafaelcalleja/go-kit/internal/common/tests/mysql_tests"
	"sync"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type TransactionInitializerExecutorSimpleDb struct {
	db       *sql.DB
	tx       *sql.Tx
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

	conn, _ := mysql_tests.NewMySQLConnection()
	i.db = conn.DB

	tx, err := i.db.Begin()
	i.tx = tx
	i.executor.Set(tx)

	return i, err
}

func (i *TransactionInitializerExecutorSimpleDb) Rollback() error {
	defer i.mu.Unlock()
	defer i.executor.Set(i.db)

	return i.tx.Rollback()
}

func (i *TransactionInitializerExecutorSimpleDb) Commit() error {
	defer i.mu.Unlock()
	defer i.executor.Set(i.db)

	return i.tx.Commit()
}

package adapters

import (
	"context"
	"database/sql"
	"sync"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type TransactionConnectionSqlInitializer struct {
	mu         sync.Mutex
	connection *ConnectionSql
	tx         *sql.Tx
}

func NewTransactionConnectionSqlInitializer(connection *ConnectionSql) *TransactionConnectionSqlInitializer {
	return &TransactionConnectionSqlInitializer{
		connection: connection,
	}
}

func (i *TransactionConnectionSqlInitializer) Begin(ctx context.Context) (transaction.Transaction, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if err := i.connection.db.Ping(); err != nil {
		return nil, err
	}

	tx, err := i.connection.db.BeginTx(ctx, nil)

	i.tx = tx

	return tx, err
}

func (i *TransactionConnectionSqlInitializer) Rollback() (err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	defer func() {
		i.tx = nil
	}()

	return i.tx.Rollback()
}

func (i *TransactionConnectionSqlInitializer) Commit() error {
	i.mu.Lock()
	defer i.mu.Unlock()
	defer func() {
		i.tx = nil
	}()

	return i.tx.Commit()
}

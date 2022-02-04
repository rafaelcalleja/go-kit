package adapters

import (
	"context"
	"sync"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type TransactionConnectionInitializer struct {
	mu         sync.Mutex
	connection *ConnectionSqlShared
}

func NewTransactionConnectionInitializer(connection *ConnectionSqlShared) *TransactionConnectionInitializer {
	return &TransactionConnectionInitializer{
		connection: connection,
	}
}

func (i *TransactionConnectionInitializer) Begin(ctx context.Context) (transaction.Transaction, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if err := i.connection.db.Ping(); err != nil {
		return nil, err
	}

	tx, err := i.connection.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	i.connection.replace(tx)

	return i, err
}

func (i *TransactionConnectionInitializer) Rollback() (err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	defer i.connection.replace(i.connection.db)

	return i.connection.get().(transaction.Transaction).Rollback()
}

func (i *TransactionConnectionInitializer) Commit() error {
	i.mu.Lock()
	defer i.mu.Unlock()
	defer i.connection.replace(i.connection.db)

	return i.connection.get().(transaction.Transaction).Commit()
}

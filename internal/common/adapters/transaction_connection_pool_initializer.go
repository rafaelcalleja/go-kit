package adapters

import (
	"context"
	"database/sql"
	"github.com/rafaelcalleja/go-kit/uuid"
	"sync"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type TransactionConnectionPoolInitializer struct {
	mu         sync.Mutex
	connection *ConnectionSql
	tx         *sql.Tx
	pool       *sync.Map
}

func NewTransactionConnectionPoolInitializer(pool *sync.Map, connection *ConnectionSql) *TransactionConnectionPoolInitializer {
	return &TransactionConnectionPoolInitializer{
		connection: connection,
		pool:       pool,
	}
}

func (i *TransactionConnectionPoolInitializer) Begin(ctx context.Context) (transaction.Transaction, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if err := i.connection.db.Ping(); err != nil {
		return nil, err
	}

	tx, err := i.connection.db.BeginTx(ctx, nil)

	txId := ctx.Value(transaction.CtxSessionIdKey.String())
	if nil == txId {
		txId = uuid.New().Create()
	}

	i.pool.Store(txId, tx)

	i.tx = tx

	return tx, err
}

func (i *TransactionConnectionPoolInitializer) Rollback() (err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	defer func() {
		i.tx = nil
	}()

	return i.tx.Rollback()
}

func (i *TransactionConnectionPoolInitializer) Commit() error {
	i.mu.Lock()
	defer i.mu.Unlock()
	defer func() {
		i.tx = nil
	}()

	return i.tx.Commit()
}

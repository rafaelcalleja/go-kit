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
	connection *ConnectionPoolSql
	tx         *sql.Tx
	pool       *sync.Map
	txId       uuid.UUID
}

func NewTransactionConnectionPoolInitializer(pool *sync.Map, connection *ConnectionPoolSql) *TransactionConnectionPoolInitializer {
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

	i.txId = txId.(uuid.UUID)
	i.pool.Store(i.txId, tx)

	i.tx = tx

	return i, err
}

func (i *TransactionConnectionPoolInitializer) Rollback() (err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	defer func() {
		i.tx = nil
		i.pool.Delete(i.txId)
	}()

	return i.tx.Rollback()
}

func (i *TransactionConnectionPoolInitializer) Commit() error {
	i.mu.Lock()
	defer i.mu.Unlock()
	defer func() {
		i.tx = nil
		i.pool.Delete(i.txId)
	}()

	return i.tx.Commit()
}

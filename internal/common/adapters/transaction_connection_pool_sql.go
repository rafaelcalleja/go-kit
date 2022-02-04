package adapters

import (
	"context"
	"database/sql"
	"sync"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type ConnectionPoolSql struct {
	connection transaction.TxPool
	db         *sql.DB
	mu         sync.Mutex
}

func (e *ConnectionPoolSql) Begin(ctx context.Context) (transaction.Transaction, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.db.Ping(); err != nil {
		return nil, err
	}

	tx, err := e.db.BeginTx(ctx, nil)
	if nil != err {
		return nil, err
	}

	txId := e.connection.StoreTransaction(ctx, tx)

	return e.connection.GetTransaction(txId), err
}

func (e *ConnectionPoolSql) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.connection.GetConnection(ctx).ExecContext(ctx, query, args...)
}

func (e *ConnectionPoolSql) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.connection.GetConnection(ctx).PrepareContext(ctx, query)
}

func (e *ConnectionPoolSql) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.connection.GetConnection(ctx).QueryContext(ctx, query, args...)
}

func (e *ConnectionPoolSql) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.connection.GetConnection(ctx).QueryRowContext(ctx, query, args...)
}

func NewConnectionPoolInitializerSql(connection transaction.TxPool) transaction.Initializer {
	conn := &ConnectionPoolSql{
		connection: connection,
		db:         connection.GetConnection(context.Background()).(*sql.DB),
	}

	return conn
}

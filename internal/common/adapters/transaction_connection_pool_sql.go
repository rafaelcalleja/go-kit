package adapters

import (
	"context"
	"database/sql"
	"sync"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type ConnectionPoolSql struct {
	*ConnectionSql
	pool *sync.Map
}

func (e *ConnectionPoolSql) resolveConnection(ctx context.Context) transaction.Connection {
	txId := ctx.Value(transaction.CtxSessionIdKey.String())

	if nil != txId {
		if cn, ok := e.pool.Load(txId); ok {
			return cn.(transaction.Connection)
		}
	}

	return e.db
}

func (e *ConnectionPoolSql) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.resolveConnection(ctx).ExecContext(ctx, query, args...)
}

func (e *ConnectionPoolSql) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.resolveConnection(ctx).PrepareContext(ctx, query)
}

func (e *ConnectionPoolSql) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.resolveConnection(ctx).QueryContext(ctx, query, args...)
}

func (e *ConnectionPoolSql) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.resolveConnection(ctx).QueryRowContext(ctx, query, args...)
}

func NewConnectionPoolSql(pool *sync.Map, db *sql.DB) transaction.Connection {
	conn := &ConnectionPoolSql{
		ConnectionSql: &ConnectionSql{
			object: db,
			db:     db,
		},
		pool: pool,
	}

	return conn
}

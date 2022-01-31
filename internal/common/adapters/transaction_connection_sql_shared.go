package adapters

import (
	"context"
	"database/sql"
	"sync"

	common_sync "github.com/rafaelcalleja/go-kit/internal/common/domain/sync"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type ConnectionSqlShared struct {
	conn   *ConnectionSql
	mu     sync.RWMutex
	locker *common_sync.ChanSync
	db     *sql.DB
}

func (e *ConnectionSqlShared) get() transaction.Connection {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.conn.object
}

func (e *ConnectionSqlShared) replace(i transaction.Connection) transaction.Connection {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.conn.object = i

	return e
}

func (e *ConnectionSqlShared) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	dispatching := ctx.Value("command_bus_dispatching")
	if true == e.locker.ChanInUse() && nil == dispatching {
		e.locker.LockAndWait()
		defer e.locker.Unlock()
	}

	return e.conn.ExecContext(ctx, query, args...)
}

func (e *ConnectionSqlShared) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	dispatching := ctx.Value("command_bus_dispatching")
	if true == e.locker.ChanInUse() && nil == dispatching {
		e.locker.LockAndWait()
		defer e.locker.Unlock()
	}

	return e.conn.PrepareContext(ctx, query)
}

func (e *ConnectionSqlShared) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	dispatching := ctx.Value("command_bus_dispatching")
	if true == e.locker.ChanInUse() && nil == dispatching {
		e.locker.LockAndWait()
		defer e.locker.Unlock()
	}

	return e.conn.QueryContext(ctx, query, args...)
}

func (e *ConnectionSqlShared) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	dispatching := ctx.Value("command_bus_dispatching")
	if true == e.locker.ChanInUse() && nil == dispatching {
		e.locker.LockAndWait()
		defer e.locker.Unlock()
	}

	return e.conn.QueryRowContext(ctx, query, args...)
}

func NewConnectionSqlShared(db *sql.DB, locker *common_sync.ChanSync) transaction.Connection {
	ex := new(ConnectionSqlShared)
	connectionSql := NewConnectionSql(db)

	ex.conn = connectionSql.(*ConnectionSql)
	ex.locker = locker
	ex.db = ex.conn.db

	return ex
}

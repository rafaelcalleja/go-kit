package adapters

import (
	"context"
	"database/sql"
	"sync"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type ConnectionSql struct {
	object transaction.Connection
	mu     sync.RWMutex
	db     *sql.DB
}

func (e *ConnectionSql) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	switch e.object.(type) {
	case *sql.Tx:
		return e.object.(*sql.Tx).ExecContext(ctx, query, args...)
	}

	return e.object.(*sql.DB).ExecContext(ctx, query, args...)
}

func (e *ConnectionSql) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	switch e.object.(type) {
	case *sql.Tx:
		return e.object.(*sql.Tx).PrepareContext(ctx, query)
	}

	return e.object.(*sql.DB).PrepareContext(ctx, query)
}

func (e *ConnectionSql) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	switch e.object.(type) {
	case *sql.Tx:
		return e.object.(*sql.Tx).QueryContext(ctx, query, args...)
	}

	return e.object.(*sql.DB).QueryContext(ctx, query, args...)
}

func (e *ConnectionSql) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	e.mu.Lock()
	defer e.mu.Unlock()

	switch e.object.(type) {
	case *sql.Tx:
		return e.object.(*sql.Tx).QueryRowContext(ctx, query, args...)
	}

	return e.object.(*sql.DB).QueryRowContext(ctx, query, args...)
}

func NewConnectionSql(db *sql.DB) transaction.Connection {
	ex := new(ConnectionSql)
	ex.object = db
	ex.db = db
	return ex
}

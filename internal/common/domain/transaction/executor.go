package transaction

import (
	"context"
	"database/sql"
	"sync"
)

type Executor interface {
	GetConnection() Connection
	WithConnection(connection Connection) Executor
	Connection
}

type ExecutorDefault struct {
	mu     sync.RWMutex
	object interface{}
}

func (e *ExecutorDefault) GetConnection() Connection {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.object.(Connection)
}

func (e *ExecutorDefault) WithConnection(i Connection) Executor {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.object = i

	return e
}

func (e *ExecutorDefault) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	switch e.object.(type) {
	case *sql.Tx:
		return e.object.(*sql.Tx).ExecContext(ctx, query, args...)
	}

	return e.object.(*sql.DB).ExecContext(ctx, query, args...)
}

func (e *ExecutorDefault) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	switch e.object.(type) {
	case *sql.Tx:
		return e.object.(*sql.Tx).PrepareContext(ctx, query)
	}

	return e.object.(*sql.DB).PrepareContext(ctx, query)
}

func (e *ExecutorDefault) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	switch e.object.(type) {
	case *sql.Tx:
		return e.object.(*sql.Tx).QueryContext(ctx, query, args...)
	}

	return e.object.(*sql.DB).QueryContext(ctx, query, args...)
}

func (e *ExecutorDefault) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	e.mu.Lock()
	defer e.mu.Unlock()
	switch e.object.(type) {
	case *sql.Tx:
		return e.object.(*sql.Tx).QueryRowContext(ctx, query, args...)
	}

	return e.object.(*sql.DB).QueryRowContext(ctx, query, args...)
}

func NewExecutor() Executor {
	ex := new(ExecutorDefault)

	return ex
}

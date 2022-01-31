package transaction

import (
	"context"
	"database/sql"
	"sync"

	common_sync "github.com/rafaelcalleja/go-kit/internal/common/domain/sync"
)

type Executor interface {
	GetConnection() Connection
	WithConnection(connection Connection) Executor
	Connection
}

type ExecutorDefault struct {
	mu     sync.RWMutex
	object interface{}
	locker *common_sync.ChanSync
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
	dispatching := ctx.Value("command_bus_dispatching")
	if true == e.locker.ChanInUse() && nil == dispatching {
		e.locker.LockAndWait()
		defer e.locker.Unlock()
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	switch e.object.(type) {
	case *sql.Tx:
		return e.object.(*sql.Tx).ExecContext(ctx, query, args...)
	}

	return e.object.(*sql.DB).ExecContext(ctx, query, args...)
}

func (e *ExecutorDefault) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	dispatching := ctx.Value("command_bus_dispatching")
	if true == e.locker.ChanInUse() && nil == dispatching {
		e.locker.LockAndWait()
		defer e.locker.Unlock()
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	switch e.object.(type) {
	case *sql.Tx:
		return e.object.(*sql.Tx).PrepareContext(ctx, query)
	}

	return e.object.(*sql.DB).PrepareContext(ctx, query)
}

func (e *ExecutorDefault) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	dispatching := ctx.Value("command_bus_dispatching")
	if true == e.locker.ChanInUse() && nil == dispatching {
		e.locker.LockAndWait()
		defer e.locker.Unlock()
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	switch e.object.(type) {
	case *sql.Tx:
		return e.object.(*sql.Tx).QueryContext(ctx, query, args...)
	}

	return e.object.(*sql.DB).QueryContext(ctx, query, args...)
}

func (e *ExecutorDefault) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	dispatching := ctx.Value("command_bus_dispatching")
	if true == e.locker.ChanInUse() && nil == dispatching {
		e.locker.LockAndWait()
		defer e.locker.Unlock()
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	switch e.object.(type) {
	case *sql.Tx:
		return e.object.(*sql.Tx).QueryRowContext(ctx, query, args...)
	}

	return e.object.(*sql.DB).QueryRowContext(ctx, query, args...)
}

func NewExecutor(locker *common_sync.ChanSync) Executor {
	ex := new(ExecutorDefault)
	ex.locker = locker

	return ex
}

package transaction

import (
	"context"
	"database/sql"
)

type Initializer interface {
	Begin(ctx context.Context) (Transaction, error)
}

type Transaction interface {
	Rollback() error
	Commit() error
}

type Connection interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

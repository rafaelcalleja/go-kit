package transaction

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrUnableToGenerateNewSession = errors.New("unable to generate new session")

	_ Querier = &sql.DB{}
	_ Querier = &sql.Tx{}
	_ Querier = &sql.Conn{}
)

type transactionKey struct{}

type Initializer interface {
	Begin(ctx context.Context) (Transaction, error)
}

type Transaction interface {
	Rollback() error
	Commit() error
}

type TxQuerier interface {
	Transaction
	Querier
}

type Querier interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type TxHandler interface {
	ManageTransaction(ctx context.Context, transaction TxQuerier) (TxId, error)
	GetTransaction(txId TxId) (Transaction, error)
}

type SafeQuerier interface {
	Get(ctx context.Context) Querier
}

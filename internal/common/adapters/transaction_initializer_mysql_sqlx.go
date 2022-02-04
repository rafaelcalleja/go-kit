package adapters

import (
	"context"
	"database/sql"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type TransactionInitializerDb struct {
	db *sql.DB
}

func NewTransactionInitializerDb(db *sql.DB) TransactionInitializerDb {
	if db == nil {
		panic("missing db")
	}

	return TransactionInitializerDb{
		db: db,
	}
}

func (i TransactionInitializerDb) Begin(ctx context.Context) (transaction.Transaction, error) {
	tx, err := i.db.BeginTx(ctx, nil)

	t := transaction.Transaction(tx)

	return t, err
}

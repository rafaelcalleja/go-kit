package adapters

import (
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

func (i TransactionInitializerDb) Begin() (transaction.Transaction, error) {
	tx, err := i.db.Begin()

	t := transaction.Transaction(tx)

	return t, err
}

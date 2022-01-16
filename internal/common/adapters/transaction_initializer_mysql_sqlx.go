package adapters

import (
	"github.com/jmoiron/sqlx"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type TransactionInitializerDb struct {
	db *sqlx.DB
}

func NewTransactionInitializerDb(db *sqlx.DB) TransactionInitializerDb {
	if db == nil {
		panic("missing db")
	}

	return TransactionInitializerDb{
		db: db,
	}
}

func (i TransactionInitializerDb) Begin() (transaction.Transaction, error) {
	tx, err := i.db.Beginx()
	wrapper := transaction.NewTransactionWrapper(tx)

	return wrapper, err
}

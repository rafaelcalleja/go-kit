package adapters

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type TransactionSessionDb struct {
	db *sqlx.DB
}

func NewTransactionSessionDb(db *sqlx.DB) TransactionSessionDb {
	if db == nil {
		panic("missing db")
	}

	return TransactionSessionDb{
		db: db,
	}
}

func (s *TransactionSessionDb) ExecuteAtomically(operation transaction.Operation) error {
	tx, err := s.db.Beginx()

	if err != nil {
		return errors.Wrap(err, "unable to start transaction")
	}

	defer func() {
		err = s.finishTransaction(err, tx)
	}()

	return operation()
}

func (s *TransactionSessionDb) finishTransaction(err error, tx *sqlx.Tx) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrap(err, "unable to rollback transaction")
		}

		return err
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			return errors.Wrap(err, "failed to commit tx")
		}

		return nil
	}
}

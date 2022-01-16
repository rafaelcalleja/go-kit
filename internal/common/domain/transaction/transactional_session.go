package transaction

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	ErrUnableToStartTransaction    = errors.New("unable to start transaction")
	ErrUnableToRollbackTransaction = errors.New("unable to rollback transaction")
	ErrUnableToCommitTransaction   = errors.New("unable to commit transaction")
	ErrPanicInOperation            = errors.New("panic in operation")
	ErrPanicInTransaction          = errors.New("panic in transaction")
)

type SessionInitializer struct {
	initializer Initializer
}

func NewTransactionalSession(initializer Initializer) SessionInitializer {
	return SessionInitializer{
		initializer: initializer,
	}
}

func (s *SessionInitializer) ExecuteAtomically(operation Operation) (err error) {
	var tx Transaction

	defer func() {
		if p := recover(); p != nil {
			switch p.(type) {
			case string:
				err = fmt.Errorf("%w: %s", ErrPanicInOperation, p.(string))
			default:
				err = ErrPanicInOperation
			}
		}

		if nil == tx {
			return
		}

		err = s.finishTransaction(err, tx.(Transaction))
	}()

	tx, err = s.initializer.Begin()

	if err != nil {
		return fmt.Errorf("%w: %s", err, ErrUnableToStartTransaction.Error())
	}

	return operation()
}

func (s *SessionInitializer) finishTransaction(err error, tx Transaction) (txErr error) {
	defer func() {
		if p := recover(); p != nil {
			switch p.(type) {
			case string:
				txErr = fmt.Errorf("%w: %s", ErrPanicInTransaction, p.(string))
			default:
				txErr = ErrPanicInTransaction
			}
		}
	}()

	if err != nil {
		txErr = err
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			txErr = fmt.Errorf("%w: %s", rollbackErr, ErrUnableToRollbackTransaction.Error())

		}
		return
	} else {
		txErr = nil
		if commitErr := tx.Commit(); commitErr != nil {
			txErr = fmt.Errorf("%w: %s", commitErr, ErrUnableToCommitTransaction.Error())

		}
		return
	}
}

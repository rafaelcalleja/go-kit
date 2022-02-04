package transaction

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/rafaelcalleja/go-kit/uuid"
)

var (
	ErrUnableToStartTransaction    = errors.New("unable to start transaction")
	ErrInitializerNilTransaction   = errors.New("initializer nil transaction")
	ErrUnableToRollbackTransaction = errors.New("unable to rollback transaction")
	ErrUnableToCommitTransaction   = errors.New("unable to commit transaction")
	ErrPanicInOperation            = errors.New("panic in operation")
	ErrPanicInTransaction          = errors.New("panic in transaction")
)

type sessionKey string

func (c sessionKey) String() string {
	return "atomic_session_" + string(c)
}

var (
	CtxSessionIdKey = sessionKey("id")
)

type SessionInitializer struct {
	initializer Initializer
	mutex       sync.Mutex
}

func NewTransactionalSession(initializer Initializer) *SessionInitializer {
	return &SessionInitializer{
		initializer: initializer,
	}
}

func (s *SessionInitializer) ExecuteAtomically(ctx context.Context, operation Operation) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	ctx = context.WithValue(ctx, CtxSessionIdKey.String(), uuid.New().String(uuid.New().Create()))

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

	tx, err = s.initializer.Begin(ctx)

	if nil == tx {
		return ErrInitializerNilTransaction
	}

	if err != nil {
		return fmt.Errorf("%w: %s", err, ErrUnableToStartTransaction.Error())
	}

	return operation(ctx)
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
	}

	txErr = nil
	if commitErr := tx.Commit(); commitErr != nil {
		txErr = fmt.Errorf("%w: %s", commitErr, ErrUnableToCommitTransaction.Error())
	}

	return
}

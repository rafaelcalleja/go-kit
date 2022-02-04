package transaction

import "context"

type Operation func(context.Context) error

type TransactionalSession interface {
	ExecuteAtomically(context.Context, Operation) error
}

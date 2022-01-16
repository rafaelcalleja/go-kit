package transaction

type Operation func() error

type TransactionalSession interface {
	ExecuteAtomically(Operation) error
}

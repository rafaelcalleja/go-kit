package transaction

type Initializer interface {
	Begin() (Transaction, error)
}

type Transaction interface {
	Rollback() error
	Commit() error
}

type Wrapper struct {
	object Transaction
}

func NewTransactionWrapper(object interface{}) Wrapper {
	return Wrapper{
		object: object.(Transaction),
	}
}

func (w Wrapper) Rollback() error {
	return w.object.Rollback()
}

func (w Wrapper) Commit() error {
	return w.object.Commit()
}

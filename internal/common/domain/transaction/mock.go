package transaction

type MockInitializer struct {
	BeginFn func() (tx Transaction, err error)
}

func NewMockInitializer() MockInitializer {
	return MockInitializer{
		BeginFn: func() (tx Transaction, err error) { return NewMockTransaction(), nil },
	}
}

func (m MockInitializer) Begin() (tx Transaction, err error) {
	return m.BeginFn()
}

type MockTransaction struct {
	RollbackFn func() error
	CommitFn   func() error
}

func NewMockTransaction() MockTransaction {
	return MockTransaction{
		RollbackFn: func() error { return nil },
		CommitFn:   func() error { return nil },
	}
}

func (m MockTransaction) Rollback() error {
	return m.RollbackFn()
}

func (m MockTransaction) Commit() error {
	return m.CommitFn()
}

type MockTransactionalSession struct {
	ExecuteAtomicallyFn func(Operation) error
}

func NewTransactionalSessionMock() MockTransactionalSession {
	return MockTransactionalSession{
		ExecuteAtomicallyFn: func(operation Operation) error { return nil },
	}
}

func (m MockTransactionalSession) ExecuteAtomically(operation Operation) error {
	return m.ExecuteAtomicallyFn(operation)
}

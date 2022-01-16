package transaction

type MockInitializer struct {
	BeginFn func() (tx Transaction, err error)
}

func NewMockInitializer() MockInitializer {
	return MockInitializer{
		BeginFn: func() (tx Transaction, err error) { return nil, nil },
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

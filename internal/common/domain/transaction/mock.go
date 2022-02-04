package transaction

import "context"

type MockInitializer struct {
	BeginFn func(ctx context.Context) (tx Transaction, err error)
}

func NewMockInitializer() MockInitializer {
	return MockInitializer{
		BeginFn: func(ctx context.Context) (tx Transaction, err error) { return NewMockTransaction(), nil },
	}
}

func (m MockInitializer) Begin(ctx context.Context) (tx Transaction, err error) {
	return m.BeginFn(ctx)
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
	ExecuteAtomicallyFn func(context.Context, Operation) error
}

func NewTransactionalSessionMock() MockTransactionalSession {
	return MockTransactionalSession{
		ExecuteAtomicallyFn: func(ctx context.Context, operation Operation) error { return nil },
	}
}

func (m MockTransactionalSession) ExecuteAtomically(ctx context.Context, operation Operation) error {
	return m.ExecuteAtomicallyFn(ctx, operation)
}

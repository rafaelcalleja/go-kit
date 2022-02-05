package transaction

import (
	"context"
	"database/sql"
)

type MockInitializer struct {
	BeginFn func(ctx context.Context) (tx Transaction, err error)
}

func NewMockInitializer() Initializer {
	return &MockInitializer{
		BeginFn: func(ctx context.Context) (tx Transaction, err error) { return NewMockTransaction(), nil },
	}
}

func (m *MockInitializer) Begin(ctx context.Context) (tx Transaction, err error) {
	return m.BeginFn(ctx)
}

type MockTransaction struct {
	RollbackFn func() error
	CommitFn   func() error
	Connection
}

func NewMockTransaction() Transaction {
	return &MockTransaction{
		RollbackFn: func() error { return nil },
		CommitFn:   func() error { return nil },
		Connection: NewMockConnection(),
	}
}

func (m *MockTransaction) Rollback() error {
	return m.RollbackFn()
}

func (m *MockTransaction) Commit() error {
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

func (m *MockTransactionalSession) ExecuteAtomically(ctx context.Context, operation Operation) error {
	return m.ExecuteAtomicallyFn(ctx, operation)
}

type MockConnection struct {
	ExecContextFn     func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContextFn  func(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContextFn    func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContextFn func(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func NewMockConnection() Connection {
	return &MockConnection{
		ExecContextFn:     func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) { return nil, nil },
		PrepareContextFn:  func(ctx context.Context, query string) (*sql.Stmt, error) { return nil, nil },
		QueryContextFn:    func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) { return nil, nil },
		QueryRowContextFn: func(ctx context.Context, query string, args ...interface{}) *sql.Row { return nil },
	}
}

func (m *MockConnection) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.ExecContextFn(ctx, query, args...)
}

func (m *MockConnection) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return m.PrepareContextFn(ctx, query)
}

func (m *MockConnection) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return m.QueryContextFn(ctx, query, args...)
}

func (m *MockConnection) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return m.QueryRowContextFn(ctx, query, args...)
}

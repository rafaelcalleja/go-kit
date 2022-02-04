package transaction

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSessionInitializer_ExecuteAtomically(t *testing.T) {
	ctx := context.Background()
	nilOperation := func(context.Context) error { return nil }
	errInOperation := errors.New("mock operation")

	errorOperation := func(context.Context) error {
		return errInOperation
	}

	panicOperation := func(context.Context) error {
		panic("mock panic")
	}

	t.Run("tx start failed", func(t *testing.T) {
		t.Parallel()
		mockInitializer := NewMockInitializer()
		mockErr := errors.New("mock tx")

		mockInitializer.BeginFn = func(_ context.Context) (tx Transaction, err error) {
			return NewMockTransaction(), mockErr
		}

		session := NewTransactionalSession(mockInitializer)
		err := session.ExecuteAtomically(ctx, nilOperation)
		require.Error(t, err)
		require.Error(t, errors.Unwrap(err))
		require.ErrorIs(t, mockErr, errors.Unwrap(err))
	})

	t.Run("tx commit failed", func(t *testing.T) {
		t.Parallel()
		mockInitializer := NewMockInitializer()
		mockTransaction := NewMockTransaction()
		mockErr := errors.New("mock commit")

		mockInitializer.BeginFn = func(_ context.Context) (tx Transaction, err error) {
			return mockTransaction, nil
		}

		mockTransaction.CommitFn = func() error {
			return mockErr
		}

		session := NewTransactionalSession(mockInitializer)
		err := session.ExecuteAtomically(ctx, nilOperation)

		require.Error(t, err)
		require.Error(t, errors.Unwrap(err))
		require.ErrorIs(t, mockErr, errors.Unwrap(err))
	})

	t.Run("tx rollback failed", func(t *testing.T) {
		t.Parallel()
		mockInitializer := NewMockInitializer()
		mockTransaction := NewMockTransaction()
		mockErr := errors.New("mock rollback")

		mockInitializer.BeginFn = func(_ context.Context) (tx Transaction, err error) {
			return mockTransaction, nil
		}

		mockTransaction.RollbackFn = func() error {
			return mockErr
		}

		session := NewTransactionalSession(mockInitializer)
		err := session.ExecuteAtomically(ctx, errorOperation)

		require.Error(t, err)
		require.Error(t, errors.Unwrap(err))
		require.ErrorIs(t, mockErr, errors.Unwrap(err))
	})

	t.Run("tx start success and committed", func(t *testing.T) {
		t.Parallel()
		mockInitializer := NewMockInitializer()
		mockTransaction := NewMockTransaction()
		mockInitializer.BeginFn = func(_ context.Context) (tx Transaction, err error) {
			return mockTransaction, nil
		}

		calledRollback := false
		mockTransaction.RollbackFn = func() error {
			calledRollback = true
			return nil
		}

		calledCommit := false
		mockTransaction.CommitFn = func() error {
			calledCommit = true
			return nil
		}

		session := NewTransactionalSession(mockInitializer)
		err := session.ExecuteAtomically(ctx, nilOperation)
		require.NoError(t, err)
		require.True(t, calledCommit)
		require.False(t, calledRollback)
	})

	t.Run("tx start success and rollbacked", func(t *testing.T) {
		t.Parallel()
		mockInitializer := NewMockInitializer()
		mockTransaction := NewMockTransaction()
		mockInitializer.BeginFn = func(_ context.Context) (tx Transaction, err error) {
			return mockTransaction, nil
		}

		calledRollback := false
		mockTransaction.RollbackFn = func() error {
			calledRollback = true
			return nil
		}

		calledCommit := false
		mockTransaction.CommitFn = func() error {
			calledCommit = true
			return nil
		}

		session := NewTransactionalSession(mockInitializer)
		err := session.ExecuteAtomically(ctx, errorOperation)
		require.Error(t, err)
		require.ErrorIs(t, errInOperation, err)

		require.False(t, calledCommit)
		require.True(t, calledRollback)
	})

	t.Run("panic in operation", func(t *testing.T) {
		t.Parallel()
		mockInitializer := NewMockInitializer()
		mockTransaction := NewMockTransaction()
		mockInitializer.BeginFn = func(_ context.Context) (tx Transaction, err error) {
			return mockTransaction, nil
		}

		calledRollback := false
		mockTransaction.RollbackFn = func() error {
			calledRollback = true
			return nil
		}

		calledCommit := false
		mockTransaction.CommitFn = func() error {
			calledCommit = true
			return nil
		}

		session := NewTransactionalSession(mockInitializer)
		err := session.ExecuteAtomically(ctx, panicOperation)
		require.Error(t, err)
		require.ErrorIs(t, ErrPanicInOperation, errors.Unwrap(err))

		require.False(t, calledCommit)
		require.True(t, calledRollback)
	})

	t.Run("unexpected content when panic in operation", func(t *testing.T) {
		t.Parallel()
		mockInitializer := NewMockInitializer()
		mockTransaction := NewMockTransaction()
		mockInitializer.BeginFn = func(_ context.Context) (tx Transaction, err error) {
			return mockTransaction, nil
		}

		calledRollback := false
		mockTransaction.RollbackFn = func() error {
			calledRollback = true
			return nil
		}

		calledCommit := false
		mockTransaction.CommitFn = func() error {
			calledCommit = true
			return nil
		}

		session := NewTransactionalSession(mockInitializer)
		err := session.ExecuteAtomically(ctx, func(_ context.Context) error {
			panic(errors.New("error in panic"))
		})

		require.Error(t, err)
		require.ErrorIs(t, ErrPanicInOperation, err)

		require.False(t, calledCommit)
		require.True(t, calledRollback)
	})

	t.Run("panic in tx start", func(t *testing.T) {
		t.Parallel()
		mockInitializer := NewMockInitializer()

		mockInitializer.BeginFn = func(_ context.Context) (tx Transaction, err error) {
			panic("panic from tx begin")
		}

		session := NewTransactionalSession(mockInitializer)
		err := session.ExecuteAtomically(ctx, nilOperation)
		require.Error(t, err)
		require.ErrorIs(t, ErrPanicInOperation, errors.Unwrap(err))
	})

	t.Run("panic in tx commit", func(t *testing.T) {
		t.Parallel()
		mockInitializer := NewMockInitializer()
		mockTransaction := NewMockTransaction()

		mockInitializer.BeginFn = func(_ context.Context) (tx Transaction, err error) {
			return mockTransaction, nil
		}

		mockTransaction.CommitFn = func() error {
			panic("panic from tx commit")
		}

		session := NewTransactionalSession(mockInitializer)
		err := session.ExecuteAtomically(ctx, nilOperation)
		require.Error(t, err)
		require.ErrorIs(t, ErrPanicInTransaction, errors.Unwrap(err))
	})

	t.Run("panic in tx rollback", func(t *testing.T) {
		t.Parallel()
		mockInitializer := NewMockInitializer()
		mockTransaction := NewMockTransaction()

		mockInitializer.BeginFn = func(_ context.Context) (tx Transaction, err error) {
			return mockTransaction, nil
		}

		mockTransaction.RollbackFn = func() error {
			panic("panic from tx rollback")
		}

		session := NewTransactionalSession(mockInitializer)
		err := session.ExecuteAtomically(ctx, errorOperation)
		require.Error(t, err)
		require.ErrorIs(t, ErrPanicInTransaction, errors.Unwrap(err))
	})

	t.Run("initializer return nil transaction without err", func(t *testing.T) {
		t.Parallel()
		mockInitializer := NewMockInitializer()

		mockInitializer.BeginFn = func(_ context.Context) (tx Transaction, err error) {
			return nil, nil
		}

		session := NewTransactionalSession(mockInitializer)
		err := session.ExecuteAtomically(ctx, errorOperation)
		require.Error(t, err)
		require.ErrorIs(t, ErrInitializerNilTransaction, err)
	})
}

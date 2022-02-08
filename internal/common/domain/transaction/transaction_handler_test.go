package transaction

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTxFromPool_GetConnection(t *testing.T) {
	t.Parallel()

	t.Run("get default connection on empty context", func(t *testing.T) {
		ctx := context.Background()

		expected := NewMockConnection()
		pool := NewTxHandler(expected)
		actual := pool.Get(ctx)

		require.Same(t, expected, actual)
	})

	t.Run("get stored connection on context with atomic session id", func(t *testing.T) {
		notExpected := NewMockConnection()
		expected := NewMockTransaction()

		pool := NewTxHandler(notExpected)
		txId, err := pool.ManageTransaction(context.Background(), expected)
		require.NoError(t, err)

		ctx := context.WithValue(context.Background(), transactionKey{}, txId.String())
		actual := pool.Get(ctx)

		require.NotSame(t, notExpected, actual)
		require.Same(t, expected, actual)
	})

	t.Run("previous atomic session exists", func(t *testing.T) {
		defaultDB := NewMockConnection()
		txDB := NewMockTransaction()

		pool := NewTxHandler(defaultDB)

		atomicSessionId := "bc6359ec-18da-420e-aa35-6a4758de04f6"
		ctx := context.WithValue(context.Background(), transactionKey{}, atomicSessionId)

		txId, err := pool.ManageTransaction(ctx, txDB)
		require.NoError(t, err)
		require.Same(t, txDB, pool.Get(ctx))
		require.Same(t, defaultDB, pool.Get(context.Background()))

		pool.repo.Delete(ctx, txId)
		require.Same(t, defaultDB, pool.Get(ctx))
	})

	t.Run("empty atomic session generate new txId", func(t *testing.T) {
		mockConnection := NewMockConnection()
		pool := NewTxHandler(mockConnection)
		mockTransaction := NewMockTransaction()
		calledCounter := 0

		commitError := errors.New("commit error")
		mockTransaction.(*MockTransaction).CommitFn = func() error {
			calledCounter++
			if calledCounter > 1 {
				return commitError
			}

			return nil
		}

		txId, err := pool.ManageTransaction(context.Background(), mockTransaction)
		require.NoError(t, err)
		transaction, err := pool.GetTransaction(txId)
		require.NoError(t, err)
		require.Same(
			t,
			mockTransaction,
			pool.Get(
				context.WithValue(context.Background(), transactionKey{}, txId.String()),
			),
		)

		require.True(t, transaction.(*txHandler).bg.id.Equals(&txId))
		require.Same(t, transaction.(*txHandler).bg.transaction, mockTransaction)
		require.Same(t, transaction.(*txHandler).repo, pool.repo)

		err = transaction.Commit()
		require.Equal(t, calledCounter, 1)
		require.NoError(t, err)
		_, err = pool.GetTransaction(txId)
		require.Error(t, ErrTransactionNotManaged, err)

		err = transaction.Commit()
		require.ErrorIs(t, commitError, err)
	})

	t.Run("store multiple transaction from multiples contexts", func(t *testing.T) {
		pool := NewTxHandler(NewMockConnection())

		_, _ = pool.ManageTransaction(context.Background(), NewMockTransaction())
		_, _ = pool.ManageTransaction(context.Background(), NewMockTransaction())
		txId, _ := pool.ManageTransaction(context.Background(), NewMockTransaction())
		_, err := pool.ManageTransaction(context.WithValue(context.Background(), transactionKey{}, txId.String()), NewMockTransaction())
		require.Error(t, ErrTransactionIdDuplicated, err)

		require.Equal(t, 3, pool.repo.Len())
	})

	t.Run("rollback remove tx", func(t *testing.T) {
		pool := NewTxHandler(NewMockConnection())
		mockTransaction := NewMockTransaction()

		calledCounter := 0
		rollbackErr := errors.New("rollback error")
		mockTransaction.(*MockTransaction).RollbackFn = func() error {
			calledCounter++
			if calledCounter > 1 {
				return rollbackErr
			}

			return nil
		}

		txId, err := pool.ManageTransaction(context.Background(), mockTransaction)
		require.NoError(t, err)
		transaction, err := pool.GetTransaction(txId)

		err = transaction.Rollback()
		require.Equal(t, calledCounter, 1)
		require.NoError(t, err)

		_, err = pool.GetTransaction(txId)
		require.Error(t, ErrTransactionNotManaged, err)

		err = transaction.Rollback()
		require.ErrorIs(t, rollbackErr, err)
	})

}

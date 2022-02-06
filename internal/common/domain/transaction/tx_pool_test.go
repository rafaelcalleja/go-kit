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
		pool := NewTxPool(expected)
		actual := pool.GetConnection(ctx)

		require.Same(t, expected, actual)
	})

	t.Run("get stored connection on context with atomic session id", func(t *testing.T) {
		notExpected := NewMockConnection()
		expected := NewMockTransaction()

		pool := NewTxPool(notExpected)
		txId, err := pool.StoreTransaction(context.Background(), expected)
		require.NoError(t, err)

		ctx := context.WithValue(context.Background(), transactionKey{}, txId.String())
		actual := pool.GetConnection(ctx)

		require.NotSame(t, notExpected, actual)
		require.Same(t, expected, actual)
	})

	t.Run("previous atomic session exists", func(t *testing.T) {
		defaultDB := NewMockConnection()
		txDB := NewMockTransaction()

		pool := NewTxPool(defaultDB)

		atomicSessionId := "bc6359ec-18da-420e-aa35-6a4758de04f6"
		ctx := context.WithValue(context.Background(), transactionKey{}, atomicSessionId)

		txId, err := pool.StoreTransaction(ctx, txDB)
		require.NoError(t, err)
		require.Same(t, txDB, pool.GetConnection(ctx))
		require.Same(t, defaultDB, pool.GetConnection(context.Background()))

		pool.RemoveTransaction(txId)
		require.Same(t, defaultDB, pool.GetConnection(ctx))
	})
	t.Run("empty atomic session generate new txId", func(t *testing.T) {
		mockConnection := NewMockConnection()
		pool := NewTxPool(mockConnection)
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

		txId, err := pool.StoreTransaction(context.Background(), mockTransaction)
		require.NoError(t, err)
		transaction, err := pool.GetTransaction(txId)
		require.NoError(t, err)
		require.Same(
			t,
			mockTransaction,
			pool.GetConnection(
				context.WithValue(context.Background(), transactionKey{}, txId.String()),
			),
		)

		require.Same(t, transaction.(*txFromPool).db, mockConnection)
		require.True(t, transaction.(*txFromPool).txId.Equals(&txId))
		require.Same(t, transaction.(*txFromPool).tx, mockTransaction)
		require.Same(t, transaction.(*txFromPool).Map, pool.(*txFromPool).Map)

		err = transaction.Commit()
		require.Equal(t, calledCounter, 1)
		require.NoError(t, err)
		_, err = pool.GetTransaction(txId)
		require.Error(t, ErrTransactionNotFound, err)

		err = transaction.Commit()
		require.ErrorIs(t, commitError, err)
	})

	t.Run("store multiple transaction from multiples contexts", func(t *testing.T) {
		pool := NewTxPool(NewMockConnection())

		_, _ = pool.StoreTransaction(context.Background(), NewMockTransaction())
		_, _ = pool.StoreTransaction(context.Background(), NewMockTransaction())
		txId, _ := pool.StoreTransaction(context.Background(), NewMockTransaction())
		_, err := pool.StoreTransaction(context.WithValue(context.Background(), transactionKey{}, txId.String()), NewMockTransaction())
		require.Error(t, ErrTransactionIdDuplicated, err)

		require.Equal(t, 3, pool.(*txFromPool).Len())
	})
}

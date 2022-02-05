package pool

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewTransactionPoolSuccess(t *testing.T) {
	ctx := context.TODO()
	mockInitializer := transaction.NewMockInitializer()
	mockTransaction := transaction.NewMockTransaction()

	commitCalledCount := 0
	mockTransaction.(*transaction.MockTransaction).CommitFn = func() error {
		commitCalledCount = commitCalledCount + 1
		return nil
	}

	rollbackCalledCount := 0
	mockTransaction.(*transaction.MockTransaction).RollbackFn = func() error {
		rollbackCalledCount = rollbackCalledCount + 1
		return nil
	}

	mockInitializer.(*transaction.MockInitializer).BeginFn = func(ctx context.Context) (tx transaction.Transaction, err error) {
		return mockTransaction, nil
	}

	mockPool := NewMockPool()

	poolGetCounter := 0
	mockPool.GetFn = func(ctx context.Context) interface{} {
		poolGetCounter = poolGetCounter + 1
		return mockInitializer
	}

	poolReleaseCounter := 0
	mockPool.ReleaseFn = func() {
		poolReleaseCounter = poolReleaseCounter + 1
	}

	txPool := NewTransactionPoolInitializer(mockPool)

	tx, _ := txPool.Begin(ctx)
	tx2, _ := txPool.Begin(ctx)
	_ = tx2.Rollback()
	_ = tx.Commit()

	require.Equal(t, 1, commitCalledCount)
	require.Equal(t, 1, rollbackCalledCount)
	require.Equal(t, 2, poolGetCounter)
	require.Equal(t, 1, poolReleaseCounter)
}

func TestNewTransactionPoolError(t *testing.T) {
	ctx := context.TODO()
	mockPool := NewMockPool()
	mockInitializer := transaction.NewMockInitializer()

	mockPool.GetFn = func(ctx context.Context) interface{} {
		return mockInitializer
	}

	txPool := NewTransactionPoolInitializer(mockPool)
	tx, _ := txPool.Begin(ctx)
	_ = tx.Commit()
	err := tx.Commit()

	require.Error(t, err)

	tx2, _ := txPool.Begin(ctx)
	_ = tx2.Rollback()
	err2 := tx.Rollback()
	require.Error(t, err2)
}

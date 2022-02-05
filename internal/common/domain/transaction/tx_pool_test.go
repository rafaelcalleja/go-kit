package transaction

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
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
		txId := pool.StoreTransaction(context.Background(), expected)

		ctx := context.WithValue(context.Background(), ctxSessionIdKey.String(), txId.String())
		actual := pool.GetConnection(ctx)

		require.NotSame(t, notExpected, actual)
		require.Same(t, expected, actual)
	})
}

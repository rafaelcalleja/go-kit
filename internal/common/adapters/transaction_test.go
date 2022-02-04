package adapters

import (
	"context"
	"testing"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
	"github.com/stretchr/testify/assert"
)

func TestNewTransactionInitializerDb(t *testing.T) {
	s := transaction.NewTransactionalSession(
		transaction.NewMockInitializer(),
	)

	ctx := context.Background()

	called := false

	_ = s.ExecuteAtomically(ctx, func(_ context.Context) error {
		called = true

		return nil
	})

	assert.True(t, called)
}

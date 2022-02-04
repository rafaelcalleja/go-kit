package adapters

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
	"github.com/rafaelcalleja/go-kit/internal/common/tests/mysql_tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTransactionInitializerDb(t *testing.T) {
	connection, _ := mysql_tests.NewMySQLConnection()

	s := transaction.NewTransactionalSession(
		NewTransactionInitializerDb(connection),
	)

	ctx := context.Background()

	called := false

	_ = s.ExecuteAtomically(ctx, func() error {
		called = true

		return nil
	})

	assert.True(t, called)
}

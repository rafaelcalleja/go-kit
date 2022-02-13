package commands

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

// TransactionalBus is an implementation of the commands.Handler.
type TransactionalBus struct {
	Bus
	session transaction.TransactionalSession
}

func NewTransactionalCommandBus(commandBus Bus, session transaction.TransactionalSession) *TransactionalBus {
	return &TransactionalBus{
		Bus:     commandBus,
		session: session,
	}
}

// Dispatch implements the commands.Bus interface.
func (b *TransactionalBus) Dispatch(ctx context.Context, cmd Command) error {
	return b.session.ExecuteAtomically(ctx, func(ctx context.Context) error {
		return b.Bus.Dispatch(ctx, cmd)
	})
}

// Register implements the commands.Bus interface.
func (b *TransactionalBus) Register(cmdType Type, handler Handler) {
	b.Bus.Register(cmdType, handler)
}

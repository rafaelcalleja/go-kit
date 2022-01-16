package commands

import (
	"context"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

// TransactionalBus is an implementation of the commands.Handler.
type TransactionalBus struct {
	commandBus CommandBus
	session    transaction.TransactionalSession
}

func NewTransactionalCommandBus(commandBus CommandBus, session transaction.TransactionalSession) TransactionalBus {
	return TransactionalBus{
		commandBus: commandBus,
		session:    session,
	}
}

// Dispatch implements the commands.Bus interface.
func (b *TransactionalBus) Dispatch(ctx context.Context, cmd Command) error {
	return b.session.ExecuteAtomically(func() error {
		return b.commandBus.Dispatch(ctx, cmd)
	})
}

// Register implements the commands.Bus interface.
func (b *TransactionalBus) Register(cmdType Type, handler Handler) {
	b.commandBus.Register(cmdType, handler)
}

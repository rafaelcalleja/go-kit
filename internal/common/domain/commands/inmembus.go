package commands

import (
	"context"
)

// CommandBus is an in-memory implementation of the command.Bus.
type CommandBus struct {
	handlers map[Type]Handler
}

// NewInMemCommandBus initializes a new instance of CommandBus.
func NewInMemCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[Type]Handler),
	}
}

// Dispatch implements the command.Bus interface.
func (b *CommandBus) Dispatch(ctx context.Context, cmd Command) error {
	handler, ok := b.handlers[cmd.Type()]
	if !ok {
		return nil
	}

	return handler.Handle(ctx, cmd)
}

// Register implements the command.Bus interface.
func (b *CommandBus) Register(cmdType Type, handler Handler) {
	b.handlers[cmdType] = handler
}

package commands

import (
	"context"
)

// CommandBus is an in-memory implementation of the commands.Bus.
type CommandBus struct {
	handlers map[Type]Handler
	pipeline *Pipeline
}

// NewInMemCommandBus initializes a new instance of CommandBus.
func NewInMemCommandBus() Bus {
	return &CommandBus{
		handlers: make(map[Type]Handler),
		pipeline: NewPipeline(),
	}
}

// Dispatch implements the commands.Bus interface.
func (b *CommandBus) Dispatch(ctx context.Context, cmd Command) error {
	handler, ok := b.handlers[cmd.Type()]
	if !ok {
		return nil
	}

	return b.handle(handler, ctx, cmd)
}

func (b *CommandBus) handle(handler Handler, ctx context.Context, cmd Command) error {
	return b.pipeline.Handle(handler, ctx, cmd)
}

func (b *CommandBus) UseMiddleware(middleware ...Middleware) {
	b.pipeline.Add(middleware...)
}

// Register implements the commands.Bus interface.
func (b *CommandBus) Register(cmdType Type, handler Handler) {
	b.handlers[cmdType] = handler
}

package commands

import (
	"context"
	"fmt"
)

// CommandBus is an in-memory implementation of the commands.Bus.
type CommandBus struct {
	handlers map[Type]Handler
	chain    MiddlewaresChain
}

// NewInMemCommandBus initializes a new instance of CommandBus.
func NewInMemCommandBus() Bus {
	return &CommandBus{
		handlers: make(map[Type]Handler),
		chain:    make([]MiddlewareFunc, 0),
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
	b.UseMiddleware(func(stack MiddlewaresChain, ctx context.Context, command Command) error {
		defer fmt.Printf("POST - execute %s\n", command.Type())
		fmt.Printf("PRE - execute %s\n", command.Type())
		return stack.Next(ctx, cmd)
	})

	newChain := append(b.chain, func(stack MiddlewaresChain, ctx context.Context, command Command) error {
		return handler.Handle(ctx, cmd)
	})

	return newChain.Next(ctx, cmd)
}

func (b *CommandBus) UseMiddleware(middleware ...MiddlewareFunc) {
	b.chain = append(b.chain, middleware...)
}

// Register implements the commands.Bus interface.
func (b *CommandBus) Register(cmdType Type, handler Handler) {
	b.handlers[cmdType] = handler
}

package commands

import (
	"context"
)

// CommandBusMock is a mock implementation of the commands.Bus.
type CommandBusMock struct {
	DispatchFn func(ctx context.Context, command Command) error
	RegisterFn func(cmdType Type, handler Handler)
}

// NewMockCommandBus initializes a new CommandBus.
func NewMockCommandBus() *CommandBusMock {
	return &CommandBusMock{
		DispatchFn: func(ctx context.Context, command Command) error { return nil },
		RegisterFn: func(cmdType Type, handler Handler) {},
	}
}

// Dispatch implements the commands.Bus interface.
func (b *CommandBusMock) Dispatch(ctx context.Context, command Command) error {
	return b.DispatchFn(ctx, command)
}

// Register implements the commands.Bus interface.
func (b *CommandBusMock) Register(cmdType Type, handler Handler) {
	b.RegisterFn(cmdType, handler)
}

type mockCommand struct{}

func newMockCommand() mockCommand {
	return mockCommand{}
}

func (command mockCommand) Type() Type {
	return "mock.command.type"
}

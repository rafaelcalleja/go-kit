package commands

import (
	"context"
)

// Bus defines the expected behaviour from a command bus.
type Bus interface {
	// Dispatch is the method used to dispatch new commands.
	Dispatch(ctx context.Context, command Command) error
	// Register is the method used to register a new command handler.
	Register(cmdType Type, handler Handler)
}

// Type represents an application command type.
type Type string

// Command represents an application command.
type Command interface {
	Type() Type
}

// Handler defines the expected behaviour from a command handler.
type Handler interface {
	Handle(context.Context, Command) error
}

package queries

import "context"

// Bus defines the expected behaviour from a query bus.
type Bus interface {
	// Ask is the method used to dispatch new queries.
	Ask(context.Context, Query) error
	// Register is the method used to register a new command handler.
	Register(Type, Handler)
}

// Type represents an application command type.
type Type string

// Query represents an application query.
type Query interface {
	Type() Type
}

// Handler defines the expected behaviour from a command handler.
type Handler interface {
	Handle(context.Context, Query) error
}

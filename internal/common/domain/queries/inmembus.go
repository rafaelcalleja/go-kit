package queries

import (
	"context"
)

// QueryBus is an in-memory implementation of the queries.Bus.
type QueryBus struct {
	handlers map[Type]Handler
}

// NewInMemQueryBus initializes a new instance of QueryBus.
func NewInMemQueryBus() *QueryBus {
	return &QueryBus{
		handlers: make(map[Type]Handler),
	}
}

// Ask implements the queries.Bus interface.
func (b *QueryBus) Ask(ctx context.Context, query Query) error {
	handler, ok := b.handlers[query.Type()]
	if !ok {
		return nil
	}

	return handler.Handle(ctx, query)
}

// Register implements the queries.Bus interface.
func (b *QueryBus) Register(cmdType Type, handler Handler) {
	b.handlers[cmdType] = handler
}

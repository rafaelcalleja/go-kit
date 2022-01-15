package queries

import (
	"context"
)

// QueryBusMock is a mock implementation of the queries.Bus.
type QueryBusMock struct {
	AskFn      func(ctx context.Context, query Query) error
	RegisterFn func(cmdType Type, handler Handler)
}

// NewMockQueryBus initializes a new QueryBus.
func NewMockQueryBus() *QueryBusMock {
	return &QueryBusMock{
		AskFn:      func(ctx context.Context, query Query) error { return nil },
		RegisterFn: func(cmdType Type, handler Handler) {},
	}
}

// Ask implements the queries.Bus interface.
func (b *QueryBusMock) Ask(ctx context.Context, query Query) error {
	return b.AskFn(ctx, query)
}

// Register implements the queries.Bus interface.
func (b *QueryBusMock) Register(cmdType Type, handler Handler) {
	b.RegisterFn(cmdType, handler)
}

package commands

import (
	"context"
	"fmt"
)

// MiddlewareFunc defines the handler used by middleware as return value.
type MiddlewareFunc func(stack MiddlewaresChain, ctx context.Context, command Command) error

// MiddlewaresChain defines a MiddlewareFunc slice.
type MiddlewaresChain []MiddlewareFunc

// Last returns the last handler in the chain. i.e. the last handler is the main one.
func (c MiddlewaresChain) Last() MiddlewareFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

func (c MiddlewaresChain) First() MiddlewareFunc {
	return c[0]
}

func (c MiddlewaresChain) Debug() {
	fmt.Printf("len=%d cap=%d %v\n", c.Len(), c.Cap(), c)
}

func (c MiddlewaresChain) Len() int {
	return len(c)
}

func (c MiddlewaresChain) Cap() int {
	return cap(c)
}

func (c *MiddlewaresChain) Next(ctx context.Context, command Command) error {
	if 0 == c.Len() {
		return nil
	}

	next := (*c)[0:c.Cap()][0]
	*c = (*c)[1:c.Cap()]

	return next(*c, ctx, command)
}

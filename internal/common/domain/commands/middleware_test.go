package commands

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMiddlewaresChain_Next(t *testing.T) {
	var chain MiddlewaresChain = make([]MiddlewareFunc, 0)

	chain = append(chain, func(stack MiddlewaresChain, ctx context.Context, command Command) error {
		defer fmt.Println("POST - first")
		fmt.Println("PRE - first")
		return stack.Next(ctx, command)
	}, func(stack MiddlewaresChain, ctx context.Context, command Command) error {
		defer fmt.Println("POST - last")
		fmt.Println("PRE - last")
		return stack.Next(ctx, command)
	})

	err := chain.Next(context.Background(), newMockCommand())
	require.NoError(t, err)

	/*

		middleware := chain.Next()
		chain.Debug()

		err := middleware(context.Background(), newMockCommand())
		require.NoError(t, err)

		chain.Next()
		chain.Debug()*/
}

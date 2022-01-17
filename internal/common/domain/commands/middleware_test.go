package commands

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/middleware"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDefaultPipeline_Handle(t *testing.T) {
	mockCommand := newMockCommand()
	mockHandler := newMockHandler()
	ctx := context.Background()

	pipeline := NewPipeline()

	countCalled := 0
	middlewareA := NewMiddlewareFunc(func(stack middleware.StackMiddleware, closure middleware.Closure, currentCtx context.Context, command Command) error {
		countCalled++
		require.Same(t, ctx, currentCtx)
		require.Equal(t, mockCommand.Type(), command.Type())
		return stack.Next().Handle(stack, closure)
	})

	calledHandlerCounter := 0
	mockHandler.HandleFn = func(currentCtx context.Context, command Command) error {
		require.Same(t, ctx, currentCtx)
		require.Equal(t, mockCommand.Type(), command.Type())
		calledHandlerCounter++
		return nil
	}

	pipeline.Add(middlewareA, middlewareA)

	err := pipeline.Handle(mockHandler, ctx, mockCommand)
	require.NoError(t, err)

	require.Equal(t, 1, calledHandlerCounter)
	require.Equal(t, 2, countCalled)

	err = pipeline.Handle(mockHandler, ctx, mockCommand)
	require.NoError(t, err)
	require.Equal(t, 2, calledHandlerCounter)
	require.Equal(t, 4, countCalled)
}

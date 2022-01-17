package commands

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/middleware"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestCommandBus_Dispatch(t *testing.T) {
	ctx := context.Background()

	mockCommand := newMockCommand()
	mockHandler := newMockHandler()

	t.Run("Pipeline is executed during command handle", func(t *testing.T) {
		called := false
		mockPipeline := newMockPipeline()
		mockPipeline.HandleFn = func(handler Handler, currentCtx context.Context, command Command) error {
			require.Same(t, reflect.TypeOf(mockHandler), reflect.TypeOf(handler))
			require.Same(t, ctx, currentCtx)
			require.Equal(t, mockCommand.Type(), command.Type())
			called = true
			return nil
		}

		commandBus, err := NewInMemCommandBusWith(
			InMemCommandBusWithPipeline(mockPipeline),
		)
		require.NoError(t, err)

		commandBus.Register(mockCommand.Type(), mockHandler)

		err = commandBus.Dispatch(ctx, mockCommand)
		require.NoError(t, err)
		require.True(t, called)
	})

	t.Run("Middlewares are added to collaborator", func(t *testing.T) {
		mockPipeline := newMockPipeline()
		called := false
		mockPipeline.AddFn = func(middlewares ...Middleware) {
			called = true
		}

		commandBus, err := NewInMemCommandBusWith(
			InMemCommandBusWithPipeline(mockPipeline),
		)
		require.NoError(t, err)

		commandBus.UseMiddleware(NewMiddlewareFunc(
			func(stack middleware.StackMiddleware, closure middleware.Closure, ctx context.Context, cmd Command) error {
				return nil
			}),
		)

		require.True(t, called)

	})
}

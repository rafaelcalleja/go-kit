package commands

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransactionalBus_Dispatch(t *testing.T) {
	sessionMock := transaction.NewTransactionalSessionMock()
	commandBusMock := NewMockCommandBus()

	ctx := context.Background()

	calledExecuteAtomically := false
	sessionMock.ExecuteAtomicallyFn = func(ctx context.Context, operation transaction.Operation) error {
		calledExecuteAtomically = true
		return operation()
	}

	calledDispatch := false
	commandBusMock.DispatchFn = func(ctx context.Context, command Command) error {
		calledDispatch = true
		return nil
	}

	commandBus := NewTransactionalCommandBus(
		commandBusMock,
		sessionMock,
	)

	err := commandBus.Dispatch(ctx, newMockCommand())
	require.NoError(t, err)
	require.True(t, calledExecuteAtomically)
	require.True(t, calledDispatch)
}

func TestTransactionalBus_Register(t *testing.T) {
	sessionMock := transaction.NewTransactionalSessionMock()
	commandBusMock := NewMockCommandBus()

	calledExecuteAtomically := false
	sessionMock.ExecuteAtomicallyFn = func(ctx context.Context, operation transaction.Operation) error {
		calledExecuteAtomically = true
		return operation()
	}

	calledDispatch := false
	commandBusMock.DispatchFn = func(ctx context.Context, command Command) error {
		calledDispatch = true
		return nil
	}

	calledRegister := false
	commandBusMock.RegisterFn = func(cmdType Type, handler Handler) {
		calledRegister = true
	}

	commandBus := NewTransactionalCommandBus(
		commandBusMock,
		sessionMock,
	)

	command := newMockCommand()
	commandBus.Register(command.Type(), nil)
	require.True(t, calledRegister)
	require.False(t, calledExecuteAtomically)
	require.False(t, calledDispatch)

}

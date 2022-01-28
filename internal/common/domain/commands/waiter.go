package commands

import (
	"context"
)

// WaiterBus is an implementation of the commands.Handler.
type WaiterBus struct {
	commandBus Bus
	wait       chan bool
}

func NewWaiterBus(commandBus Bus, wait chan bool) *WaiterBus {
	return &WaiterBus{
		commandBus: commandBus,
		wait:       wait,
	}
}

// Dispatch implements the commands.Bus interface.
func (b *WaiterBus) Dispatch(ctx context.Context, cmd Command) error {
	<-b.wait
	defer func() {
		b.wait <- true
	}()
	return b.commandBus.Dispatch(ctx, cmd)
}

// Register implements the commands.Bus interface.
func (b *WaiterBus) Register(cmdType Type, handler Handler) {
	b.commandBus.Register(cmdType, handler)
}

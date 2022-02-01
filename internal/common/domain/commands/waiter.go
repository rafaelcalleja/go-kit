package commands

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/sync"
)

// WaiterBus is an implementation of the commands.Handler.
type WaiterBus struct {
	commandBus Bus
	locker     *sync.ChanSync
}

func NewWaiterBus(commandBus Bus, locker *sync.ChanSync) *WaiterBus {
	return &WaiterBus{
		commandBus: commandBus,
		locker:     locker,
	}
}

// Dispatch implements the commands.Bus interface.
func (b *WaiterBus) Dispatch(ctx context.Context, cmd Command) error {
	ctx = b.locker.Lock(ctx)
	defer b.locker.Unlock()

	return b.commandBus.Dispatch(ctx, cmd)
}

// Register implements the commands.Bus interface.
func (b *WaiterBus) Register(cmdType Type, handler Handler) {
	b.commandBus.Register(cmdType, handler)
}

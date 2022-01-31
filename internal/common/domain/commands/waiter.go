package commands

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/pool"
	"sync"
)

// WaiterBus is an implementation of the commands.Handler.
type WaiterBus struct {
	commandBus Bus
	wait       chan bool
	mu         *sync.RWMutex
	sem        *pool.Semaphore
	cond       *sync.Cond
}

func NewWaiterBus(commandBus Bus, wait chan bool, mu *sync.RWMutex, sem *pool.Semaphore, cond *sync.Cond) *WaiterBus {
	return &WaiterBus{
		commandBus: commandBus,
		wait:       wait,
		mu:         mu,
		sem:        sem,
		cond:       cond,
	}
}

// Dispatch implements the commands.Bus interface.
func (b *WaiterBus) Dispatch(ctx context.Context, cmd Command) error {
	<-b.wait
	b.cond.L.Lock()
	ctx = context.WithValue(ctx, "intx", true)
	defer func() {
		b.wait <- true
		b.cond.Broadcast()
		b.cond.L.Unlock()
	}()
	return b.commandBus.Dispatch(ctx, cmd)
}

// Register implements the commands.Bus interface.
func (b *WaiterBus) Register(cmdType Type, handler Handler) {
	b.commandBus.Register(cmdType, handler)
}

package sync

import (
	"context"
	"sync"
)

type contextMuxKey string

func (c contextMuxKey) String() string {
	return "context_mux_key_" + string(c)
}

var (
	ctxMuxKey = syncKey("locking")
)

type ContextMux struct {
	mux         Mutex
	id          int
	internalMux sync.Mutex
}

func (c *ContextMux) Lock(ctx context.Context) {
	c.internalMux.Lock()
	defer c.internalMux.Unlock()

}

func (c *ContextMux) TryLock(ctx context.Context) bool {
	c.internalMux.Lock()
	defer c.internalMux.Unlock()

	return false

}

func (c *ContextMux) IsLocked(ctx context.Context) bool {
	c.internalMux.Lock()
	defer c.internalMux.Unlock()
	return false
}

func (c *ContextMux) Unlock(ctx context.Context) {
	c.internalMux.Lock()
	defer c.internalMux.Unlock()

}

func NewContextMux() *ContextMux {
	return new(ContextMux)
}

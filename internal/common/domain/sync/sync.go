package sync

import (
	"context"
	base "sync"
)

type syncKey string

func (c syncKey) String() string {
	return "sync_key_" + string(c)
}

var (
	ctxSyncKey = syncKey("locking")
)

type ChanSync struct {
	wait chan bool
	cond *base.Cond
}

func NewChanSync() *ChanSync {
	sync := new(ChanSync)
	sync.wait = make(chan bool, 1)
	sync.cond = base.NewCond(&base.Mutex{})

	sync.wait <- true

	return sync
}

func (c *ChanSync) ChanInUse() bool {
	return len(c.wait) == 0
}

func (c *ChanSync) Lock() {
	<-c.wait
	c.cond.L.Lock()
}

func (c *ChanSync) Unlock() {
	c.wait <- true
	c.cond.L.Unlock()
	c.cond.Broadcast()
}

func (c *ChanSync) LockAndWait() {
	c.cond.L.Lock()
	for true == c.ChanInUse() {
		c.cond.Wait()
	}

	<-c.wait
}

func (c *ChanSync) CLock(ctx context.Context) context.Context {
	defer c.Lock()
	return context.WithValue(ctx, ctxSyncKey.String(), true)
}

func (c *ChanSync) CWait(ctx context.Context) bool {
	locked := ctx.Value(ctxSyncKey.String())
	if true == c.ChanInUse() && nil == locked {
		c.LockAndWait()
		return true
	}

	return false
}

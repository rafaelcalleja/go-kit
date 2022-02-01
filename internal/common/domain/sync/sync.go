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
	wait  chan bool
	cWait chan bool
	cond  *base.Cond
}

func NewChanSync() *ChanSync {
	sync := new(ChanSync)
	sync.wait = make(chan bool, 1)
	sync.cWait = make(chan bool, 1)
	sync.cond = base.NewCond(&base.Mutex{})

	sync.wait <- true

	return sync
}

func (c *ChanSync) chanInUse() bool {
	return len(c.wait) == 0
}

func (c *ChanSync) lock() {
	<-c.wait
	c.cond.L.Lock()
}

func (c *ChanSync) Unlock() {
	c.wait <- true
	c.cond.L.Unlock()
	c.cond.Broadcast()
}

func (c *ChanSync) lockAndWait() {
	c.cond.L.Lock()
	for true == c.chanInUse() {
		c.cond.Wait()
	}

	<-c.wait
	c.cWait <- true
}

func (c *ChanSync) Lock(ctx context.Context) context.Context {
	defer c.lock()
	return context.WithValue(ctx, ctxSyncKey.String(), true)
}

func (c *ChanSync) CWait(ctx context.Context) {
	locked := ctx.Value(ctxSyncKey.String())
	if true == c.chanInUse() && nil == locked {
		c.lockAndWait()
	}
}

func (c *ChanSync) CUnlock() {
	if len(c.cWait) == 1 {
		c.Unlock()
		<-c.cWait
	}
}

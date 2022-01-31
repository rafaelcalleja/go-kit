package sync

import (
	base "sync"
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

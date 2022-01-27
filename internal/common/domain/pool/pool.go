package pool

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/petermattis/goid"
	"golang.org/x/sync/semaphore"
)

type Pool interface {
	Get(ctx context.Context) interface{}
	Release()
}

type MockPool struct {
	GetFn     func(ctx context.Context) interface{}
	ReleaseFn func()
}

func NewMockPool() MockPool {
	return MockPool{
		GetFn:     func(ctx context.Context) interface{} { return nil },
		ReleaseFn: func() {},
	}
}

func (m MockPool) Get(ctx context.Context) interface{} {
	return m.GetFn(ctx)
}

func (m MockPool) Release() {
	m.ReleaseFn()
}

type SinglePool struct {
	mu     sync.Mutex
	sem    *semaphore.Weighted
	owner  int64
	object interface{}
}

func NewSingleObject(fn func() interface{}) SinglePool {
	maxWorkers := 1

	return SinglePool{
		object: fn(),
		sem:    semaphore.NewWeighted(int64(maxWorkers)),
	}
}

func (p *SinglePool) Get(ctx context.Context) interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()

	var owner int64
	atomic.StoreInt64(&owner, goid.Get())

	if p.owner == owner {
		return p.object
	}

	if err := p.sem.Acquire(ctx, 1); err != nil {
		panic("Failed to acquire semaphore")
	}

	p.owner = owner

	return p.object
}

func (p *SinglePool) Release() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.owner != goid.Get() {
		panic("Failed to release semaphore")
	}

	p.owner = -1

	p.sem.Release(1)
}

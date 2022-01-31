package pool

import (
	"context"
	"golang.org/x/sync/semaphore"
	"sync"
)

type Semaphore struct {
	mu     sync.RWMutex
	sem    *semaphore.Weighted
	object interface{}
}

type SemContext struct {
	sem    *semaphore.Weighted
	object interface{}
	mu     *sync.RWMutex
}

func NewSemaphore(fn func() interface{}) *Semaphore {
	maxWorkers := 1

	return &Semaphore{
		object: fn(),
		sem:    semaphore.NewWeighted(int64(maxWorkers)),
	}
}

func (p *Semaphore) Get(ctx context.Context) *SemContext {
	if object := ctx.Value("semaphore"); object != nil {
		return &SemContext{
			sem:    p.sem,
			object: object,
			mu:     &p.mu,
		}
	}

	if err := p.sem.Acquire(ctx, 1); err != nil {
		panic("Failed to acquire semaphore")
	}

	semContext := SemContext{
		sem:    p.sem,
		object: p.object,
		mu:     &p.mu,
	}

	return &semContext
}

func (s *SemContext) Object() interface{} {
	return s.object
}

func (s *SemContext) Lock() {
	s.mu.RLock()
}

func (s *SemContext) Unlock() {
	s.mu.RUnlock()
}

func (s *SemContext) Release() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sem.Release(1)
}

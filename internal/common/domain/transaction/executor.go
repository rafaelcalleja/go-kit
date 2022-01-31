package transaction

import (
	"sync"
)

type Executor interface {
	Get() interface{}
	Set(interface{}) Executor
}

type ExecutorDefault struct {
	mu     sync.RWMutex
	object interface{}
}

func (e *ExecutorDefault) Get() interface{} {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.object
}

func (e *ExecutorDefault) Set(i interface{}) Executor {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.object = i

	return e
}

func NewExecutor() Executor {
	ex := new(ExecutorDefault)

	return ex
}

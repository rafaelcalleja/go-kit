package transaction

import (
	"fmt"
	"github.com/petermattis/goid"
	"sync"
)

type Executor interface {
	Get() interface{}
	Set(interface{}) Executor
}

type ExecutorDefault struct {
	mu     sync.RWMutex
	object interface{}
	maps   map[int64]interface{}
}

func (e *ExecutorDefault) Get() interface{} {
	e.mu.Lock()
	defer e.mu.Unlock()

	//return e.maps[goid.Get()]
	fmt.Println("getting----------- ", goid.Get())
	return e.object
}

func (e *ExecutorDefault) Set(i interface{}) Executor {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.object = i
	e.maps[goid.Get()] = i
	fmt.Println("setting----------- ", goid.Get())
	return e
}

func NewExecutor() Executor {
	ex := new(ExecutorDefault)
	ex.maps = make(map[int64]interface{})
	return ex
}

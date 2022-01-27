package transaction

type Executor interface {
	Get() interface{}
	Set(interface{}) Executor
}

type ExecutorDefault struct {
	object interface{}
}

func (e *ExecutorDefault) Get() interface{} {
	return e.object
}

func (e *ExecutorDefault) Set(i interface{}) Executor {
	e.object = i
	return e
}

func NewExecutor() Executor {
	ex := new(ExecutorDefault)

	return ex
}

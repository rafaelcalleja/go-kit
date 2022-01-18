package middleware

type Context interface {
	Get(key string) interface{}
	Set(key string, value interface{})
	Invoke() error
	Closure() Closure
}

type EmptyCtx struct {
	values  map[string]interface{}
	closure Closure
}

func (e EmptyCtx) Get(key string) interface{} {
	val, ok := e.values[key]
	if false == ok {
		return nil
	}

	return val
}

func (e EmptyCtx) Set(key string, value interface{}) {
	(&e).values[key] = value
}

func (e EmptyCtx) Invoke() error {
	return e.closure()
}

func (e EmptyCtx) Closure() Closure {
	return e.closure
}

func ContextWith(closure Closure) EmptyCtx {
	ctx := new(EmptyCtx)
	ctx.values = make(map[string]interface{})
	ctx.closure = closure

	return *ctx
}

package middleware

type Func func(stack StackMiddleware, closure Closure) error

type Wrapper struct {
	middlewareFn Func
}

func NewMiddlewareFunc(fn Func) *Wrapper {
	return &Wrapper{
		fn,
	}
}

func (w Wrapper) Handle(stack StackMiddleware, closure Closure) error {
	return w.middlewareFn(stack, closure)
}

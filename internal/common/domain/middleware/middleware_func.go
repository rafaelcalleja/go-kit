package middleware

type Func func(stack StackMiddleware, ctx Context) error

type Wrapper struct {
	middlewareFn Func
}

func NewMiddlewareFunc(fn Func) *Wrapper {
	return &Wrapper{
		fn,
	}
}

func (w Wrapper) Handle(stack StackMiddleware, ctx Context) error {
	return w.middlewareFn(stack, ctx)
}

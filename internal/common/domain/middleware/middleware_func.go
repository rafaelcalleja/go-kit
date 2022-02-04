package middleware

import (
	"context"
)

type Func func(stack StackMiddleware, ctx context.Context, closure Closure) error

type Wrapper struct {
	middlewareFn Func
}

func NewMiddlewareFunc(fn Func) *Wrapper {
	return &Wrapper{
		fn,
	}
}

func (w Wrapper) Handle(stack StackMiddleware, ctx context.Context, closure Closure) error {
	return w.middlewareFn(stack, ctx, closure)
}

package commands

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/middleware"
)

type Middleware interface {
	Handle(stack middleware.StackMiddleware, closure middleware.Closure, ctx context.Context, cmd Command) error
}

type Func func(stack middleware.StackMiddleware, closure middleware.Closure, ctx context.Context, cmd Command) error

type Wrapper struct {
	middlewareFn Func
}

func NewMiddlewareFunc(fn Func) *Wrapper {
	return &Wrapper{
		fn,
	}
}

func (w Wrapper) Handle(stack middleware.StackMiddleware, closure middleware.Closure, ctx context.Context, cmd Command) error {
	return w.middlewareFn(stack, closure, ctx, cmd)
}

type Pipeline struct {
	middlewares []Middleware
}

func NewPipeline() *Pipeline {
	return &Pipeline{make([]Middleware, 0)}
}

func (p *Pipeline) Add(middlewares ...Middleware) {
	p.middlewares = append(p.middlewares, middlewares...)
}

func (p Pipeline) Handle(handler Handler, ctx context.Context, cmd Command) error {
	pipeline := middleware.NewPipeline()

	closure := func() error {
		return handler.Handle(ctx, cmd)
	}

	elements := make([]middleware.Middleware, len(p.middlewares))
	for i, v := range p.middlewares {
		elements[i] = middleware.NewMiddlewareFunc(func(stack middleware.StackMiddleware, closure middleware.Closure) error {
			return stack.Next().Handle(stack, func() error {
				return v.Handle(stack, closure, ctx, cmd)
			})
		})
	}

	pipeline.Add(elements...)

	return pipeline.Handle(closure)
}

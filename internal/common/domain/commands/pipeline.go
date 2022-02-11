package commands

import (
	"context"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/middleware"
)

type Pipeline interface {
	Add(middlewares ...middleware.Middleware)
	Handle(ctx context.Context, handler Handler, cmd Command) error
}

type DefaultPipeline struct {
	middlewares []middleware.Middleware
}

func NewPipeline() *DefaultPipeline {
	return &DefaultPipeline{make([]middleware.Middleware, 0)}
}

func (p *DefaultPipeline) Add(middlewares ...middleware.Middleware) {
	p.middlewares = append(p.middlewares, middlewares...)
}

func (p DefaultPipeline) Handle(ctx context.Context, handler Handler, cmd Command) error {
	pipeline := middleware.NewPipeline()

	closure := func(ctx context.Context) error {
		return handler.Handle(ctx, cmd)
	}

	pipeline.Add(middleware.NewMiddlewareFunc(func(stack middleware.StackMiddleware, ctx context.Context, closure middleware.Closure) error {
		ctx = withPipelineContext(ctx, handler, cmd)

		return stack.Next().Handle(stack, ctx, closure)
	}))

	pipeline.Add(p.middlewares...)

	return pipeline.Handle(ctx, closure)
}

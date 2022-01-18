package commands

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/middleware"
)

type Pipeline interface {
	Add(middlewares ...middleware.Middleware)
	Handle(handler Handler, ctx context.Context, cmd Command) error
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

func (p DefaultPipeline) Handle(handler Handler, ctx context.Context, cmd Command) error {
	pipeline := middleware.NewPipeline()

	closure := func() error {
		return handler.Handle(ctx, cmd)
	}

	pipeline.Add(middleware.NewMiddlewareFunc(func(stack middleware.StackMiddleware, middlewareCtx middleware.Context) error {
		middlewareCtx.Set("ctx", ctx)
		middlewareCtx.Set("command", cmd)
		middlewareCtx.Set("handler", handler)

		return stack.Next().Handle(stack, middlewareCtx)
	}))

	pipeline.Add(p.middlewares...)

	return pipeline.Handle(middleware.ContextWith(closure))
}

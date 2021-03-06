package middleware

import "context"

type Middleware interface {
	Handle(stack StackMiddleware, ctx context.Context, closure Closure) error
}

// Closure defines the handler used by middleware as return value.
type Closure func(context.Context) error

type Pipeline struct {
	stack *DefaultStackMiddleware
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		stack: NewDefaultStackMiddleware(),
	}
}

func (p *Pipeline) Add(middlewares ...Middleware) {
	elements := make([]Element, len(middlewares))
	for i, v := range middlewares {
		elements[i] = v
	}

	p.stack.Push(elements...)
}

func (p Pipeline) Handle(ctx context.Context, closure Closure) error {
	clone := p.stack.Clone()

	return clone.Next().Handle(&clone, ctx, closure)
}

package middleware

import "context"

type Middleware interface {
	Handle(stack StackMiddleware, ctx Context) error
}

// Closure defines the handler used by middleware as return value.
type Closure func() error

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

func (p Pipeline) Handle(ctx Context) error {
	clone := p.stack.Clone()

	return clone.Next().Handle(&clone, ctx)
}

const ctxDefaultContext string = "pipeline_context"

type DefaultContext struct {
	Ctx context.Context
}

func GetDefaultContext(context Context) DefaultContext {
	return context.Get(ctxDefaultContext).(DefaultContext)
}

func setDefaultContext(ctx context.Context, context Context) {
	pipelineContext := DefaultContext{
		ctx,
	}

	context.Set(ctxDefaultContext, pipelineContext)
}

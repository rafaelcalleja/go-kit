package middleware

type Middleware interface {
	Handle(stack StackMiddleware, closure Closure) error
}

// Closure defines the handler used by middleware as return value.
type Closure func() error

type Pipeline struct {
	stack *DefaultStackMiddleware
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		stack: &DefaultStackMiddleware{NewStack()},
	}
}

func (p *Pipeline) Add(middlewares ...Middleware) {
	elements := make([]Element, len(middlewares))
	for i, v := range middlewares {
		elements[i] = v
	}

	p.stack.Push(elements...)
}

func (p Pipeline) Handle(closure Closure) error {
	clone := p.stack.Clone()

	return clone.Next().Handle(&clone, closure)
}

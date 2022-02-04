package middleware

import (
	"context"
)

type StackMiddleware interface {
	Next() Middleware
}

type DefaultStackMiddleware struct {
	stack *Stack
}

func NewDefaultStackMiddlewareWith(options ...func(*DefaultStackMiddleware) error) (DefaultStackMiddleware, error) {
	var stackMiddleware = new(DefaultStackMiddleware)

	for _, option := range options {
		err := option(stackMiddleware)
		if err != nil {
			return DefaultStackMiddleware{}, err
		}
	}

	return *stackMiddleware, nil
}

func DefaultStackMiddlewareWithStack(stack *Stack) func(*DefaultStackMiddleware) error {
	return func(s *DefaultStackMiddleware) error {
		s.stack = stack
		return nil
	}
}

func NewDefaultStackMiddleware() *DefaultStackMiddleware {
	stack, _ := NewDefaultStackMiddlewareWith(
		DefaultStackMiddlewareWithStack(NewStack()),
	)

	return &stack
}

func (s *DefaultStackMiddleware) Stack() *Stack {
	return s.stack
}

func (s *DefaultStackMiddleware) Handle(_ StackMiddleware, ctx context.Context, closure Closure) error {
	return closure(ctx)
}

func (s *DefaultStackMiddleware) Next() Middleware {
	if 0 == s.stack.Size() {
		return s
	}

	return s.stack.Pop().(Middleware)
}

func (s *DefaultStackMiddleware) Push(element ...Element) {
	s.stack.Push(element...)
}

func (s DefaultStackMiddleware) Size() int {
	return s.stack.Size()
}

func (s *DefaultStackMiddleware) Pop() interface{} {
	return s.stack.Pop()
}

func (s *DefaultStackMiddleware) Clone() DefaultStackMiddleware {
	clone := s.stack.Clone()

	stackMiddleware, _ := NewDefaultStackMiddlewareWith(
		DefaultStackMiddlewareWithStack(&clone),
	)

	return stackMiddleware
}

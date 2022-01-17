package middleware

type StackMiddleware interface {
	Next() Middleware
}

type DefaultStackMiddleware struct {
	stack *Stack
}

func (s *DefaultStackMiddleware) Stack() *Stack {
	return s.stack
}

func (s *DefaultStackMiddleware) Handle(_ StackMiddleware, closure Closure) error {
	return closure()
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

	return DefaultStackMiddleware{&clone}
}

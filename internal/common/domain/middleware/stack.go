package middleware

type Stack []Element
type Element interface{}

func (s *Stack) Push(element ...Element) {
	*s = append(*s, element...)
}

func (s Stack) Size() int {
	return len(s)
}

func (s *Stack) Pop() interface{} {
	h := *s

	if 0 == s.Size() {
		return nil
	}

	var element interface{}

	element, *s = h[0], h[1:]

	return element
}

func (s Stack) Clone() Stack {
	clone := make([]Element, s.Size())
	_ = copy(clone, s)

	stack := NewStack()
	stack.Push(clone...)

	return *stack
}

func NewStack() *Stack {
	return &Stack{}
}

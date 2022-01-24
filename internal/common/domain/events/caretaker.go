package events

type CareTaker []*memento

func (s *CareTaker) Add(m ...*memento) {
	*s = append(*s, m...)
}

func (s *CareTaker) Pop() *memento {
	h := *s

	if 0 == len(h) {
		return nil
	}

	n := len(h) - 1

	element := h[n]
	*s = h[:n]

	return element
}

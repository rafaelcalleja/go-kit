package events

type memento struct {
	data []interface{}
}

func (m *memento) getData() []interface{} {
	return m.data
}

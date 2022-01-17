package middleware

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFiFoStack(t *testing.T) {
	stack := NewStack()
	stack.Push("foo")
	stack.Push("bar")

	require.Equal(t, 2, stack.Size())
	require.Equal(t, "foo", stack.Pop())
	require.Equal(t, 1, stack.Size())
	require.Equal(t, "bar", stack.Pop())
	require.Equal(t, 0, stack.Size())
}

func TestCloneStack(t *testing.T) {
	stack := NewStack()
	stack.Push("foo")
	stack.Push("bar")

	clone := *stack

	stack.Pop()
	stack.Pop()

	require.Equal(t, 2, clone.Size())
	require.Equal(t, 0, stack.Size())

	require.Equal(t, "foo", clone.Pop())
}

func TestVariadicStack(t *testing.T) {
	stack := NewStack()
	stack.Push("foo", "bar")
	require.Equal(t, 2, stack.Size())

	factory := func(text ...Element) *Stack {
		s1 := NewStack()
		s1.Push(text...)

		return s1
	}

	s2 := factory("1", "2", "3")
	require.Equal(t, "1", s2.Pop())
	require.Equal(t, 2, s2.Size())

	factoryCaster := func(text ...string) *Stack {
		elements := make([]Element, len(text))
		for i, v := range text {
			elements[i] = v
		}

		s1 := NewStack()
		s1.Push(elements...)

		return s1
	}

	s3 := factoryCaster("4", "3", "2", "1")
	require.Equal(t, "4", s3.Pop())
	require.Equal(t, 3, s3.Size())
}

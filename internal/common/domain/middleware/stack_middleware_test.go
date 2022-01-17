package middleware

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDefaultStackMiddleware_Clone(t *testing.T) {
	stack := NewDefaultStackMiddleware()
	elements := []int{1, 2, 3, 4, 5, 6}

	for _, element := range elements {
		stack.Push(element)
	}

	clone := stack.Clone()

	require.Equal(t, stack.Size(), clone.Size())
	require.NotSame(t, stack.Stack(), clone.Stack())

	clone.Pop()
	require.NotEqual(t, stack.Size(), clone.Size())
}

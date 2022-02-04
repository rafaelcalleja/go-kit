package middleware

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type printer struct {
	text   string
	logger *[]string
}

func newPrinter(text string, logger *[]string) *printer {
	return &printer{text, logger}
}

func (p *printer) Handle(stack StackMiddleware, ctx context.Context, closure Closure) error {
	*p.logger = append(*p.logger, "PRE "+p.text)

	defer func() {
		*p.logger = append(*p.logger, "POST "+p.text)
	}()

	return stack.Next().Handle(stack, ctx, closure)
}

func TestChainPipeline(t *testing.T) {
	var logger = make([]string, 0)

	pipeline := NewPipeline()
	pipeline.Add(newPrinter("printer 1", &logger), newPrinter("printer 2", &logger))

	err := pipeline.Handle(context.Background(), func(_ context.Context) error {
		logger = append(logger, "executing")
		return nil
	})

	expected := []string{"PRE printer 1", "PRE printer 2", "executing", "POST printer 2", "POST printer 1"}

	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(expected, logger))

	err = pipeline.Handle(context.Background(), func(_ context.Context) error {
		logger = append(logger, "executing 2")
		return nil
	})

	require.NoError(t, err)
	expected2 := append(expected, "PRE printer 1", "PRE printer 2", "executing 2", "POST printer 2", "POST printer 1")
	require.True(t, reflect.DeepEqual(expected2, logger))
}

func TestEmptyChainPipeline(t *testing.T) {
	pipeline := NewPipeline()

	called := false
	err := pipeline.Handle(context.Background(), func(_ context.Context) error {
		called = true
		return nil
	})

	require.NoError(t, err)
	require.True(t, called)
}

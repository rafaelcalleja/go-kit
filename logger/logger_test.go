package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFactory(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	logger := New()

	assert.Same(t, New(), logger)
}

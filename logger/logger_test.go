package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFactory(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	var logger Logger = New()
	var emptyLogger Logger = NewNullLogger()

	assert.Same(t, New(), logger)
	assert.Same(t, NewNullLogger(), emptyLogger)
	assert.NotSame(t, logger, emptyLogger)
}

package helper

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHelpOverWriteBehavior(t *testing.T) {
	err := errors.New("Dummy Error")

	called := false
	helper := NewErrorHelper()

	helper.BehaviorOnFatal(func(message string, code int) {
		assert.Contains(t, message, err.Error())
		called = true
	})

	helper.CheckErr(err)
	assert.True(t, called)
	assert.Same(t, helper, NewErrorHelper())
}

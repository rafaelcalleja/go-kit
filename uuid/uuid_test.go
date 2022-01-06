package uuid

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParse(t *testing.T) {
	stringUuid := "f11fea47-8493-415f-8e2e-401b64a1bfde"

	parse, err := New().Parse(stringUuid)
	require.NoError(t, err)

	assert.Equal(t, stringUuid, New().String(parse))
}

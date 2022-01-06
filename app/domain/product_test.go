package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEquality(t *testing.T) {
	productIdA, _ := NewProductId("fe52596a-e4fc-4887-9eee-2da9fc7d9e30")
	productIdB, _ := NewProductId("fe52596a-e4fc-4887-9eee-2da9fc7d9e30")
	productIdC, _ := NewProductId("fe52596a-e4fc-4887-9eee-2da9fc7d9e31")

	assert.True(t, productIdA.Equals(productIdB))
	assert.NotSame(t, productIdA, productIdB)

	assert.False(t, productIdA.Equals(productIdC))
}

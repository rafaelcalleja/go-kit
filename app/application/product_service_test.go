package application

import (
	"github.com/rafaelcalleja/go-kit/app/infrastructure"
	"github.com/rafaelcalleja/go-kit/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateProduct(t *testing.T) {
	t.Setenv("LOG_LEVEL", "DEBUG")
	defaultLogger := logger.NewNullLogger()
	repository := infrastructure.NewProductRepository(defaultLogger)
	service := NewProductService(repository)
	err := service.CreateProduct("1b93d80c-16b3-4338-805c-67a071db988f")

	assert.Nil(t, err)
}

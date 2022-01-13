package command

import (
	"errors"
	"github.com/rafaelcalleja/go-kit/internal/store/domain"
	"github.com/rafaelcalleja/go-kit/internal/store/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateProduct(t *testing.T) {
	repository := mock.NewMockProductRepository()

	service := NewCreateProductHandler(repository)

	t.Run("product wrong id", func(t *testing.T) {
		repository.OfFn = func(id *domain.ProductId) (*domain.Product, error) {
			return domain.NewProduct(id.String())
		}

		err := service.Handle("1")
		assert.True(t, errors.Is(err, domain.ErrWrongUuid))
	})

	t.Run("product already exists", func(t *testing.T) {
		repository.OfFn = func(id *domain.ProductId) (*domain.Product, error) {
			return domain.NewProduct(id.String())
		}

		err := service.Handle("1b93d80c-16b3-4338-805c-67a071db988f")
		assert.True(t, errors.Is(err, ErrProductAlreadyExists))
	})

	t.Run("product save successfully", func(t *testing.T) {
		repository.OfFn = func(id *domain.ProductId) (*domain.Product, error) {
			return &domain.Product{}, errors.New("product not found")
		}

		var saved string
		repository.SaveFn = func(p *domain.Product) error {
			saved = p.ID().String()
			return nil
		}

		newUuid := "1b93d80c-16b3-4338-805c-67a071db988f"
		err := service.Handle(newUuid)
		assert.Equal(t, newUuid, saved)
		assert.Nil(t, err)
	})

	t.Run("product cant be saved", func(t *testing.T) {
		repository.OfFn = func(id *domain.ProductId) (*domain.Product, error) {
			return &domain.Product{}, errors.New("product not found")
		}

		called := false
		repository.SaveFn = func(p *domain.Product) error {
			called = true
			return errors.New("error saving product")
		}

		newUuid := "1b93d80c-16b3-4338-805c-67a071db988f"
		err := service.Handle(newUuid)

		assert.True(t, called)
		assert.NotNil(t, err)
	})

}

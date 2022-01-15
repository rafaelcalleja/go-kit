package command

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/events"
	"github.com/rafaelcalleja/go-kit/internal/store/domain"
	"github.com/rafaelcalleja/go-kit/internal/store/mock"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	repository := mock.NewMockProductRepository()
	eventBus := events.NewMockEventBus()
	ctx := context.Background()

	service := NewCreateProductHandler(repository, eventBus)

	t.Run("product wrong id", func(t *testing.T) {
		repository.OfFn = func(id *domain.ProductId) (*domain.Product, error) {
			return domain.NewProduct(id.String())
		}

		err := service.Handle(ctx, NewCreateProductCommand("1"))
		assert.True(t, errors.Is(err, domain.ErrWrongUuid))
	})

	t.Run("product already exists", func(t *testing.T) {
		repository.OfFn = func(id *domain.ProductId) (*domain.Product, error) {
			return domain.NewProduct(id.String())
		}

		err := service.Handle(ctx, NewCreateProductCommand("1b93d80c-16b3-4338-805c-67a071db988f"))
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
		err := service.Handle(ctx, NewCreateProductCommand(newUuid))
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
		err := service.Handle(ctx, NewCreateProductCommand(newUuid))

		assert.True(t, called)
		assert.NotNil(t, err)
	})

	t.Run("events are published", func(t *testing.T) {
		repository.OfFn = func(id *domain.ProductId) (*domain.Product, error) {
			return &domain.Product{}, errors.New("product not found")
		}

		repository.SaveFn = func(p *domain.Product) error { return nil }

		called := false
		var collectEvents []events.Event
		eventBus.PublishFn = func(ctx context.Context, events []events.Event) error {
			called = true
			collectEvents = events
			return nil
		}

		newUuid := "1b93d80c-16b3-4338-805c-67a071db988f"
		err := service.Handle(ctx, NewCreateProductCommand(newUuid))

		assert.Nil(t, err)
		assert.True(t, called)
		assert.GreaterOrEqual(t, len(collectEvents), 1)
		assert.Equal(t, newUuid, collectEvents[0].AggregateID())
	})

	t.Run("wrong command handled", func(t *testing.T) {
		err := service.Handle(ctx, newWrongCommand())
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrProductCommandUnhandled)
	})

}

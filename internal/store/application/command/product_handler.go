package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/commands"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/events"
	"github.com/rafaelcalleja/go-kit/internal/store/domain"
)

var (
	ErrProductAlreadyExists    = errors.New("product exists")
	ErrProductCommandUnhandled = errors.New("unexpected command")
)

type CreateProductHandler struct {
	repository domain.ProductRepository
	eventBus   events.Bus
}

func NewCreateProductHandler(repository domain.ProductRepository, eventBus events.Bus) CreateProductHandler {
	return CreateProductHandler{repository, eventBus}
}

func (service CreateProductHandler) Handle(ctx context.Context, command commands.Command) error {
	createProductCommand, ok := command.(CreateProductCommand)
	if !ok {
		return ErrProductCommandUnhandled
	}

	id := createProductCommand.id
	product, err := domain.NewProduct(id)

	if nil != err {
		return err
	}

	_, err = service.repository.Of(ctx, product.ID())

	if err == nil {
		return fmt.Errorf("%w: %s", ErrProductAlreadyExists, id)
	}

	if err := service.repository.Save(ctx, product); err != nil {
		return err
	}

	return service.eventBus.Publish(ctx, product.PullEvents())
}

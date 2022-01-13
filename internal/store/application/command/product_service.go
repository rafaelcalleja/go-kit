package command

import (
	"errors"
	"fmt"

	"github.com/rafaelcalleja/go-kit/internal/store/domain"
)

var (
	ErrProductAlreadyExists = errors.New("product exists")
)

type CreateProductHandler struct {
	repository domain.ProductRepository
}

func NewCreateProductHandler(repository domain.ProductRepository) CreateProductHandler {
	return CreateProductHandler{repository}
}

func (service CreateProductHandler) Handle(id string) error {
	product, err := domain.NewProduct(id)

	if nil != err {
		return err
	}

	_, err = service.repository.Of(product.ID())

	if err == nil {
		return fmt.Errorf("%w: %s", ErrProductAlreadyExists, id)
	}

	return service.repository.Save(product)
}

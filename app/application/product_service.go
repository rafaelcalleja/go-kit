package application

import (
	"errors"
	"fmt"

	"github.com/rafaelcalleja/go-kit/app/domain"
)

var (
	ErrProductAlreadyExists = errors.New("product exists")
)

type ProductService struct {
	repository domain.ProductRepository
}

func NewProductService(repository domain.ProductRepository) *ProductService {
	return &ProductService{repository}
}

func (service ProductService) CreateProduct(id string) error {
	product, err := domain.NewProduct(id)

	if nil != err {
		return err
	}

	_, err = service.repository.Of(product.ID())

	if nil != err {
		return fmt.Errorf("%w: %s", ErrProductAlreadyExists, id)
	}

	if nil != err {
		return err
	}

	return service.repository.Save(product)
}

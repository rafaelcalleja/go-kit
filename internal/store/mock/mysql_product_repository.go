package mock

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/store/domain"
)

type ProductRepository struct {
	SaveFn func(p *domain.Product) error
	OfFn   func(id *domain.ProductId) (*domain.Product, error)
}

func NewMockProductRepository() *ProductRepository {
	return &ProductRepository{
		SaveFn: func(p *domain.Product) error { return nil },
		OfFn:   func(id *domain.ProductId) (*domain.Product, error) { return &domain.Product{}, nil },
	}
}

func (m ProductRepository) Save(_ context.Context, p *domain.Product) error {
	return m.SaveFn(p)
}

func (m ProductRepository) Of(_ context.Context, id *domain.ProductId) (*domain.Product, error) {
	return m.OfFn(id)
}

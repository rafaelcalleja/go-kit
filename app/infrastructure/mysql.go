package infrastructure

import (
	"github.com/rafaelcalleja/go-kit/app/domain"
	"github.com/rafaelcalleja/go-kit/logger"
)

type ProductRepository struct {
	logger logger.Logger
}

func NewProductRepository(logger logger.Logger) *ProductRepository {
	return &ProductRepository{logger}
}

func (m ProductRepository) Save(p *domain.Product) error {
	m.logger.Debugf("Save %s", p.ID())
	return nil
}

func (m ProductRepository) Of(id *domain.ProductId) (*domain.Product, error) {
	m.logger.Debugf("Of " + id.String())
	return domain.NewProduct(id.String())
}

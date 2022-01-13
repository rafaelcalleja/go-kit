package adapters

import (
	"github.com/rafaelcalleja/go-kit/internal/store/domain"
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
	m.logger.Debugf("Of %s", id)
	return domain.NewProduct(id.String())
}

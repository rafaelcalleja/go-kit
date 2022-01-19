package domain

import (
	"context"
	"errors"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductRepository interface {
	Save(ctx context.Context, product *Product) error
	Of(ctx context.Context, id *ProductId) (*Product, error)
}

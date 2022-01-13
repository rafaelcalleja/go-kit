package domain

import (
	"errors"
)

var (
	ErrWrongUuid = errors.New("wrong uuid")
)

type Product struct {
	id *ProductId
}

func (p Product) ID() *ProductId {
	return p.id
}

func NewProduct(id string) (*Product, error) {
	idVO, err := NewProductId(id)

	if nil != err {
		return &Product{}, err
	}

	return &Product{idVO}, nil
}

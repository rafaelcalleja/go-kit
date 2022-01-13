package domain

import (
	"errors"
	"fmt"
	"github.com/rafaelcalleja/go-kit/uuid"
)

var (
	ErrWrongUuid = errors.New("wrong uuid")
)

type ProductRepository interface {
	Save(*Product) error
	Of(id *ProductId) (*Product, error)
}

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

type ProductId struct {
	id uuid.UUID
}

func NewProductId(id string) (*ProductId, error) {
	idVO, err := uuid.New().Parse(id)

	if nil != err {
		return &ProductId{}, fmt.Errorf("%w: %s", ErrWrongUuid, id)
	}

	return &ProductId{idVO}, nil
}

func (pi *ProductId) String() string {
	return uuid.New().String(pi.id)
}

func (pi *ProductId) Equals(other *ProductId) bool {
	return other.String() == pi.String()
}

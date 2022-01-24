package domain

import (
	"errors"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/events"
)

var (
	ErrWrongUuid = errors.New("wrong uuid")
)

type Product struct {
	id *ProductId
	events.ImplementRecordableEvents
}

func (p Product) ID() *ProductId {
	return p.id
}

func NewProduct(id string) (*Product, error) {
	idVO, err := NewProductId(id)

	if nil != err {
		return &Product{}, err
	}

	product := &Product{
		id: idVO,
	}

	product.Record(NewProductCreatedEvent(idVO.String()))

	return product, nil
}

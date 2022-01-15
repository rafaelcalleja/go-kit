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

	events []events.Event
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

func (p *Product) Record(event events.Event) {
	p.events = append(p.events, event)
}

func (p *Product) PullEvents() []events.Event {
	evt := p.events
	p.events = []events.Event{}

	return evt
}

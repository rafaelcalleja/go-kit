package domain

import (
	"fmt"
	"github.com/rafaelcalleja/go-kit/uuid"
)

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

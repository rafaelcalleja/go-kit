package transaction

import (
	"errors"
	"fmt"

	"github.com/rafaelcalleja/go-kit/uuid"
)

var (
	ErrWrongUuid = errors.New("wrong uuid")
)

type TxId struct {
	id uuid.UUID
}

func NewTxId(id string) (*TxId, error) {
	idVO, err := uuid.New().Parse(id)

	if nil != err {
		return &TxId{}, fmt.Errorf("%w: %s", ErrWrongUuid, id)
	}

	return &TxId{idVO}, nil
}

func NewRandomTxId() (*TxId, error) {
	return NewTxId(uuid.New().String(uuid.New().Create()))
}

func (tx *TxId) String() string {
	return uuid.New().String(tx.id)
}

func (tx *TxId) Equals(other *TxId) bool {
	return other.String() == tx.String()
}

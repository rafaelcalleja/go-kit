package transaction

import (
	"context"
)

type Repository interface {
	Save(_ context.Context, transaction *begunTransaction) error
	Of(_ context.Context, txId TxId) (*begunTransaction, error)
	Delete(_ context.Context, txId TxId)
	Len() int
}

type begunTransaction struct {
	transaction Transaction
	id          *TxId
}

func begunTransactionWith(options ...func(*begunTransaction) error) (*begunTransaction, error) {
	var bt = new(begunTransaction)

	for _, option := range options {
		err := option(bt)

		if err != nil {
			return &begunTransaction{}, err
		}
	}

	if nil == bt.id {
		idVo, _ := NewRandomTxId()
		bt.id = idVo
	}

	return bt, nil
}

func begunTransactionWithTransaction(transaction Transaction) func(*begunTransaction) error {
	return func(s *begunTransaction) error {
		s.transaction = transaction
		return nil
	}
}

func begunTransactionWithTxId(id string) func(*begunTransaction) error {
	return func(s *begunTransaction) error {
		idVo, err := NewTxId(id)

		if nil != err {
			return err
		}

		s.id = idVo
		return nil
	}
}

func newBegunTransaction(id string, transaction Transaction) (*begunTransaction, error) {
	return begunTransactionWith(
		begunTransactionWithTransaction(transaction),
		begunTransactionWithTxId(id),
	)
}

func newBegunTransactionFromContext(ctx context.Context, transaction Transaction) (*begunTransaction, error) {
	sessionId, err := sessionIdFromContext(ctx)

	if err == nil {
		return newBegunTransaction(sessionId.String(), transaction)
	}

	return begunTransactionWith(
		begunTransactionWithTransaction(transaction),
	)
}

type begunTxRepository struct {
	connMap *connectionMap
}

func (t *begunTxRepository) Save(_ context.Context, transaction *begunTransaction) error {
	id := transaction.id.String()

	if _, exists := t.connMap.Load(id); exists == true {
		return ErrTransactionIdDuplicated
	}

	t.connMap.Store(id, transaction.transaction)

	return nil
}

func (t *begunTxRepository) Of(_ context.Context, txId TxId) (*begunTransaction, error) {
	transaction, ok := t.connMap.Load(txId.String())

	if false == ok {
		return &begunTransaction{}, ErrTransactionNotManaged
	}

	return &begunTransaction{
		transaction: transaction.(Transaction),
		id:          &txId,
	}, nil
}

func (t *begunTxRepository) Delete(_ context.Context, txId TxId) {
	t.connMap.Delete(txId.String())
}

func (t *begunTxRepository) Len() int {
	return t.connMap.Len()
}

func NewTransactionRepository(connMap *connectionMap) Repository {
	return &begunTxRepository{
		connMap: connMap,
	}
}

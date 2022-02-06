package transaction

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrTransactionNotFound     = errors.New("transaction not found")
	ErrTransactionIdDuplicated = errors.New("transaction id duplicated")
)

type txFromPool struct {
	*sync.Map
	db   Connection
	tx   Transaction
	txId TxId
	mux  sync.Mutex
}

func (t *txFromPool) GetConnection(ctx context.Context) Connection {
	txId := ctx.Value(transactionKey{})

	if nil == txId {
		return t.db
	}

	conn, _ := t.Load(txId.(string))

	switch conn.(type) {
	case Connection:
		return conn.(Connection)
	default:
		return t.db
	}
}

func NewTxPool(db Connection) TxPool {
	t := &txFromPool{
		Map: &sync.Map{},
		db:  db,
	}

	return t
}

func (t *txFromPool) GetTransaction(txId TxId) (Transaction, error) {
	transaction, ok := t.Load(txId.String())

	if false == ok {
		return &txFromPool{}, ErrTransactionNotFound
	}

	var t3 Transaction

	t3 = &txFromPool{
		Map:  t.Map,
		db:   t.db,
		tx:   transaction.(Transaction),
		txId: txId,
	}

	return t3, nil
}

func (t *txFromPool) StoreTransaction(ctx context.Context, transaction Transaction) (TxId, error) {
	txId := ctx.Value(transactionKey{})

	if nil == txId {
		idVo, _ := NewRandomTxId()
		txId = idVo.String()
	}

	txVO, err := NewTxId(txId.(string))

	if nil != err {
		return TxId{}, err
	}

	if _, exists := t.Load(txVO.String()); exists == true {
		return TxId{}, ErrTransactionIdDuplicated
	}

	t.Store(txVO.String(), transaction)

	return *txVO, nil
}

func (t *txFromPool) RemoveTransaction(txId TxId) {
	t.Delete(txId.String())
}

func (t *txFromPool) Rollback() (err error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	defer func() {
		t.RemoveTransaction(t.txId)
	}()

	return t.tx.Rollback()
}

func (t *txFromPool) Commit() error {
	t.mux.Lock()
	defer t.mux.Unlock()

	defer func() {
		t.RemoveTransaction(t.txId)
	}()

	return t.tx.Commit()
}

func (t *txFromPool) Len() int {
	length := 0
	t.Range(func(_, _ interface{}) bool {
		length++
		return true
	})

	return length
}

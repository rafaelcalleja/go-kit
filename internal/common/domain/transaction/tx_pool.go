package transaction

import (
	"context"
	"sync"
)

type txFromPool struct {
	*sync.Map
	db   Connection
	tx   Transaction
	txId TxId
	mux  sync.Mutex
}

func (t *txFromPool) GetConnection(ctx context.Context) Connection {
	txId := ctx.Value(ctxSessionIdKey.String())

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

func (t *txFromPool) GetTransaction(txId TxId) Transaction {
	transaction, _ := t.Load(txId.String())

	var t3 Transaction

	t3 = &txFromPool{
		Map:  t.Map,
		db:   t.db,
		tx:   transaction.(Transaction),
		txId: txId,
	}

	return t3
}

func (t *txFromPool) StoreTransaction(ctx context.Context, transaction Transaction) TxId {
	txId := ctx.Value(ctxSessionIdKey.String())

	if nil == txId {
		idVo, _ := NewRandomTxId()
		txId = idVo.String()
	}

	txVO, _ := NewTxId(txId.(string))

	t.Store(txVO.String(), transaction)

	return *txVO
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

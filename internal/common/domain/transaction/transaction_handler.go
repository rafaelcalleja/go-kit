package transaction

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrTransactionNotManaged   = errors.New("transaction not managed")
	ErrTransactionIdDuplicated = errors.New("transaction id duplicated")
)

type txHandler struct {
	repo Repository
	bg   *begunTransaction
	mux  sync.Mutex
}

func (t *txHandler) Get(ctx context.Context) Querier {
	return t.repo.(*begunTxRepository).connMap.LoadFromCtxOrBase(ctx)
}

func NewTxHandler(db Querier) *txHandler {
	connMap := newConnectionMap(db, &sync.Map{})

	return &txHandler{
		repo: NewTransactionRepository(
			connMap,
		),
	}
}

func (t *txHandler) SafeQuerier() SafeQuerier {
	return SafeQuerier(t)
}

func (t *txHandler) GetTransaction(txId TxId) (Transaction, error) {
	txFromRepo, err := t.repo.Of(context.Background(), txId)

	if nil != err {
		return &txHandler{}, ErrTransactionNotManaged
	}

	return &txHandler{
		repo: t.repo,
		bg:   txFromRepo,
	}, nil
}

func (t *txHandler) ManageTransaction(ctx context.Context, transaction TxQuerier) (TxId, error) {
	bt, err := newBegunTransactionFromContext(ctx, transaction)

	if nil != err {
		return TxId{}, err
	}

	if err = t.repo.Save(ctx, bt); err != nil {
		return TxId{}, ErrTransactionIdDuplicated
	}

	return *bt.id, nil
}

func (t *txHandler) Rollback() (err error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	defer func() {
		t.repo.Delete(context.Background(), *t.bg.id)
	}()

	return t.bg.transaction.Rollback()
}

func (t *txHandler) Commit() error {
	t.mux.Lock()
	defer t.mux.Unlock()

	defer func() {
		t.repo.Delete(context.Background(), *t.bg.id)
	}()

	return t.bg.transaction.Commit()
}

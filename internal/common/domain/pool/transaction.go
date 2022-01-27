package pool

import (
	"context"
	"errors"
	"sync"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

var (
	ErrTransactionNotFound = errors.New("tx not found")
)

type TransactionPool struct {
	pool    Pool
	txCount int
	mu      sync.Mutex
	ctx     context.Context
	tx      map[int]transaction.Transaction
}

func NewTransactionPoolInitializer(ctx context.Context, pool Pool) transaction.Initializer {
	return &TransactionPool{
		pool: pool,
		ctx:  ctx,
		tx:   make(map[int]transaction.Transaction),
	}
}

func (t *TransactionPool) endTransaction() {
	t.txCount--

	err := recover()
	if nil != err {
		t.txCount = 0
	}

	if 0 == t.txCount {
		t.tx = make(map[int]transaction.Transaction)
		t.pool.Release()
	}

	t.mu.Unlock()
}

func (t *TransactionPool) Rollback() (err error) {
	t.mu.Lock()
	defer t.endTransaction()

	tx, ok := t.tx[t.txCount]
	if false == ok {
		return ErrTransactionNotFound
	}

	return tx.Rollback()
}

func (t *TransactionPool) Commit() (err error) {
	t.mu.Lock()
	defer t.endTransaction()

	tx, ok := t.tx[t.txCount]
	if false == ok {
		return ErrTransactionNotFound
	}

	return tx.Commit()
}

func (t *TransactionPool) Begin() (transaction.Transaction, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	tx, err := t.pool.Get(t.ctx).(transaction.Initializer).Begin()
	t.txCount++
	t.tx[t.txCount] = tx

	return transaction.Transaction(t), err
}

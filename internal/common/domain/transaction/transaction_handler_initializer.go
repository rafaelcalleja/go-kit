package transaction

import (
	"context"
	"sync"
)

type txHandlerInitializer struct {
	initializer Initializer
	handler     TxHandler
	mu          sync.Mutex
}

func (e *txHandlerInitializer) Begin(ctx context.Context) (Transaction, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	tx, err := e.initializer.Begin(ctx)

	if nil != err {
		return nil, err
	}

	txId, err := e.handler.ManageTransaction(ctx, tx.(TxQuerier))

	if nil != err {
		return nil, err
	}

	return e.handler.GetTransaction(txId)
}

func NewTxHandlerInitializer(handler TxHandler, initializer Initializer) Initializer {
	conn := &txHandlerInitializer{
		handler:     handler,
		initializer: initializer,
	}

	return conn
}

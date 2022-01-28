package transaction

type ChainTxInitializer struct {
	chain []Initializer
}

func NewChainTxInitializer(initializer ...Initializer) *ChainTxInitializer {
	return &ChainTxInitializer{
		chain: initializer,
	}
}

func (c ChainTxInitializer) Begin() (Transaction, error) {
	txs := make([]Transaction, len(c.chain))

	for k, initializer := range c.chain {
		t, err := initializer.Begin()

		if err != nil {
			return &ChainTx{}, nil
		}

		txs[k] = t
	}

	return NewChainTx(txs), nil
}

type ChainTx struct {
	chain []Transaction
}

func NewChainTx(chain []Transaction) *ChainTx {
	return &ChainTx{
		chain: chain,
	}
}

func (c ChainTx) Rollback() (err error) {
	var err2 error
	for _, tx := range c.chain {
		err2 = tx.Rollback()

		if nil == err && err2 != nil {
			err = err2
		}
	}

	return err
}

func (c ChainTx) Commit() (err error) {
	for _, tx := range c.chain {
		err = tx.Commit()

		if nil != err {
			return err
		}
	}

	return
}

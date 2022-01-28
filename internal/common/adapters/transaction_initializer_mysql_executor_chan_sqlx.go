package adapters

import (
	"database/sql"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type TransactionInitializerExecutorChanDb struct {
	executor transaction.Executor
	conn     *sql.DB
	wait     chan bool
}

func NewTransactionInitializerExecutorChanDb(connection *sql.DB, wait chan bool, executor transaction.Executor) *TransactionInitializerExecutorChanDb {
	return &TransactionInitializerExecutorChanDb{
		executor: executor,
		conn:     connection,
		wait:     wait,
	}
}

func (i *TransactionInitializerExecutorChanDb) Begin() (transaction.Transaction, error) {
	<-i.wait

	tx, err := i.conn.Begin()

	if err != nil {
		panic(err)
	}

	i.executor.Set(tx)

	return i, err
}

func (i *TransactionInitializerExecutorChanDb) Rollback() (err error) {
	defer i.executor.Set(nil)
	defer func() {
		i.wait <- true
	}()

	err = i.executor.Get().(*sql.Tx).Rollback()

	if nil != err {

		panic(err)
	}

	return err
}

func (i *TransactionInitializerExecutorChanDb) Commit() error {
	defer i.executor.Set(nil)
	defer func() {
		i.wait <- true
	}()

	return i.executor.Get().(*sql.Tx).Commit()
}

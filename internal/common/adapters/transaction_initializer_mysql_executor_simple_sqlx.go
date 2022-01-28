package adapters

import (
	"database/sql"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
	"github.com/rafaelcalleja/go-kit/internal/common/tests/mysql_tests"
	"sync"
)

var c = 0

type TransactionInitializerExecutorSimpleDb struct {
	db       *sql.DB
	tx       *sql.Tx
	executor transaction.Executor
	mu       sync.Mutex
	txs      map[int64]*sql.Tx
	conn     *sql.DB
}

func NewTransactionInitializerExecutorSimpleDb(db *sql.DB, executor transaction.Executor) *TransactionInitializerExecutorSimpleDb {
	if db == nil {
		panic("missing db")
	}

	conn, _ := mysql_tests.NewMySQLConnection()

	return &TransactionInitializerExecutorSimpleDb{
		db:       db,
		executor: executor,
		txs:      make(map[int64]*sql.Tx),
		conn:     conn,
	}
}

func (i *TransactionInitializerExecutorSimpleDb) Begin() (transaction.Transaction, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	err := i.conn.Ping()
	if err != nil {
		panic(err)
	}

	tx, err := i.conn.Begin()
	if err != nil {
		panic(err)
	}

	i.tx = tx
	i.executor.Set(tx)

	return i, err
}

func (i *TransactionInitializerExecutorSimpleDb) Rollback() (err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	defer i.executor.Set(i.conn)

	/*delay := time.Now().Unix()/1000 - (goid.Get() * 1000)

	time.Sleep(time.Duration(delay) * time.Nanosecond)
	fmt.Printf("Delay %d of %d", int(goid.Get()), delay)*/
	switch i.executor.Get().(type) {
	case *sql.DB:
		return nil
	case *sql.Tx:
		err = i.executor.Get().(*sql.Tx).Rollback()
	}

	if nil != err {

		panic(err)
	}

	return err
}

func (i *TransactionInitializerExecutorSimpleDb) Commit() error {
	i.mu.Lock()
	defer i.mu.Unlock()
	defer i.executor.Set(i.conn)

	return i.tx.Commit()
}

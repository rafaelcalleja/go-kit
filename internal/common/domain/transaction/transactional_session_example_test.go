package transaction_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	common_adapters "github.com/rafaelcalleja/go-kit/internal/common/adapters"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

func Example_cleanArchitectureTx() {
	connection, _sqlmock, err := sqlmock.New()

	if nil != err {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50000*time.Microsecond)
	defer cancel()

	transactionHandler := transaction.NewTxHandler(connection)
	querier := transactionHandler.SafeQuerier()

	txSession := transaction.NewTransactionalSession(
		transaction.NewTxHandlerInitializer(transactionHandler, common_adapters.NewSqlDBInitializer(connection)),
	)

	service := handler{
		conn: querier,
	}

	wg := sync.WaitGroup{}

	for x := 0; x < 5; x++ {
		_sqlmock.ExpectBegin()
		_sqlmock.ExpectQuery("SELECT .*").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(4))
		_sqlmock.ExpectQuery("SELECT .*").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(6))
		_sqlmock.ExpectCommit()

		wg.Add(1)

		go func() {
			_ = txSession.ExecuteAtomically(ctx, func(ctx context.Context) error {

				service.handle(ctx, 2)
				service.handle(ctx, 3)

				wg.Done()

				return nil
			})
		}()
	}

	wg.Wait()

	for x := 0; x < 5; x++ {
		_sqlmock.ExpectQuery("SELECT .*").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(2))
		wg.Add(1)

		go func() {
			service.handle(ctx, 1)
			wg.Done()
		}()
	}

	wg.Wait()

	// Output: 4 6 4 6 4 6 4 6 4 6 2 2 2 2 2
}

type handler struct {
	conn transaction.SafeQuerier
}

func (s *handler) handle(ctx context.Context, value int) {
	rows, err := s.conn.Get(ctx).QueryContext(ctx, "SELECT (?+?) as value", value, value)

	if nil != err {
		panic(err)
	}

	defer func() {
		_ = rows.Close()
	}()

	var v float64

	for rows.Next() {

		if err = rows.Scan(&v); nil != err {
			panic(err)
		}

		fmt.Printf(" %d", int(v))
		return
	}

	panic("no rows")
}

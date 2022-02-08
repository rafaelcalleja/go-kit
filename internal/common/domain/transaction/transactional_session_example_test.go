package transaction_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	common_adapters "github.com/rafaelcalleja/go-kit/internal/common/adapters"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

var (
	connectionPointers = make(map[string]int)
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

	domainService := &handler{
		conn: querier,
	}

	eventPublisher := &handler{
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

				domainService.handle(ctx, 2)
				eventPublisher.handle(ctx, 3)

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
			domainService.handle(ctx, 1)
			wg.Done()
		}()
	}

	// [connection.(type)=sql.result connection.address]
	wg.Wait()

	// Output: [*sql.Tx=4 0x00000000]
	//[*sql.Tx=6 0x00000000]
	//[*sql.Tx=4 0x11111111]
	//[*sql.Tx=6 0x11111111]
	//[*sql.Tx=4 0x22222222]
	//[*sql.Tx=6 0x22222222]
	//[*sql.Tx=4 0x33333333]
	//[*sql.Tx=6 0x33333333]
	//[*sql.Tx=4 0x44444444]
	//[*sql.Tx=6 0x44444444]
	//[*sql.DB=2 0x55555555]
	//[*sql.DB=2 0x55555555]
	//[*sql.DB=2 0x55555555]
	//[*sql.DB=2 0x55555555]
	//[*sql.DB=2 0x55555555]
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

		//connection memory address
		c := fmt.Sprintf("%p", s.conn.Get(ctx))

		//fake connection memory address
		var id string
		if val, exists := connectionPointers[c]; exists == false {
			connectionPointers[c] = len(connectionPointers)
			id = "0x" + strings.Repeat(strconv.Itoa(connectionPointers[c]), 8)
		} else {
			id = "0x" + strings.Repeat(strconv.Itoa(val), 8)
		}

		fmt.Printf("[%T=%d %s]\n", s.conn.Get(ctx), int(v), id)
		return
	}

	panic("no rows")
}

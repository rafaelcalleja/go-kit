package adapters_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rafaelcalleja/go-kit/internal/common/adapters"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	fakeMemoryAddress = make(map[string]string)
)

func Example_cleanArchitectureTx() {
	mysqlConn, _sqlmock, err := sqlmock.New()

	if nil != err {
		panic(err)
	}

	_gorm, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mysqlConn,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if nil != err {
		panic(err)
	}

	gormDB := &adapters.GormDB{
		DB: _gorm,
	}

	if nil != err {
		panic(err)
	}

	ctx := context.Background()
	transactionHandler := transaction.NewTxHandler(gormDB)
	querier := adapters.NewGormSafeQuerier(transactionHandler.SafeQuerier())

	txSession := transaction.NewTransactionalSession(
		transaction.NewTxHandlerInitializer(transactionHandler, adapters.NewGormInitializer(gormDB)),
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
		_sqlmock.ExpectQuery("SELECT 2\\+2").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(4))
		_sqlmock.ExpectQuery("SELECT 3\\+3").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(6))
		_sqlmock.ExpectCommit()

		wg.Add(1)

		go func() {
			defer wg.Done()
			_ = txSession.ExecuteAtomically(ctx, func(ctx context.Context) error {
				domainService.handle(ctx, 2)
				eventPublisher.handle(ctx, 3)

				return nil
			})
		}()
	}

	wg.Wait()

	for x := 0; x < 5; x++ {
		_sqlmock.ExpectQuery("SELECT 1\\+1").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(2))
		wg.Add(1)

		go func() {
			defer wg.Done()
			domainService.handle(ctx, 1)
		}()
	}

	// [connection.(type)=sql.result connection.address]
	wg.Wait()

	// Output: [*gorm.DB *sql.Tx=4 0x00000000]
	//[*gorm.DB *sql.Tx=6 0x00000000]
	//[*gorm.DB *sql.Tx=4 0x11111111]
	//[*gorm.DB *sql.Tx=6 0x11111111]
	//[*gorm.DB *sql.Tx=4 0x22222222]
	//[*gorm.DB *sql.Tx=6 0x22222222]
	//[*gorm.DB *sql.Tx=4 0x33333333]
	//[*gorm.DB *sql.Tx=6 0x33333333]
	//[*gorm.DB *sql.Tx=4 0x44444444]
	//[*gorm.DB *sql.Tx=6 0x44444444]
	//[*gorm.DB *sql.DB=2 0x55555555]
	//[*gorm.DB *sql.DB=2 0x55555555]
	//[*gorm.DB *sql.DB=2 0x55555555]
	//[*gorm.DB *sql.DB=2 0x55555555]
	//[*gorm.DB *sql.DB=2 0x55555555]
}

type handler struct {
	conn *adapters.GormSafeQuerier
}

func (s *handler) handle(ctx context.Context, value int) {
	tx := s.conn.GetGorm(ctx)
	connection := s.conn.GetConnPool(ctx)

	rows, err := connection.QueryContext(ctx, fmt.Sprintf("SELECT %d+%d", value, value))

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
		c := fmt.Sprintf("%p", tx)

		//fake connection memory address
		if _, exists := fakeMemoryAddress[c]; exists == false {
			connectionCounter := len(fakeMemoryAddress)
			strCounter := strconv.Itoa(connectionCounter)
			fakeAddress := "0x" + strings.Repeat(strCounter, 8)
			fakeMemoryAddress[c] = fakeAddress
		}

		fmt.Printf("[%T %T=%d %s]\n", tx, connection, int(v), fakeMemoryAddress[c])
		return
	}

	panic("no rows")
}

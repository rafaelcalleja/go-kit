package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
	"github.com/rafaelcalleja/go-kit/internal/store/domain"
)

const (
	sqlProductTable = "products"
)

type sqlProduct struct {
	ID string `db:"id"`
}

type ProductRepository struct {
	executor  transaction.Executor
	dbTimeout time.Duration
}

func NewMysqlProductRepository(executor transaction.Executor, dbTimeout time.Duration) *ProductRepository {
	return &ProductRepository{
		executor:  executor,
		dbTimeout: dbTimeout,
	}
}

func (m ProductRepository) Save(ctx context.Context, product *domain.Product) error {
	productSQLStruct := sqlbuilder.NewStruct(new(sqlProduct))
	query, args := productSQLStruct.InsertInto(sqlProductTable, sqlProduct{
		ID: product.ID().String(),
	}).Build()

	ctxTimeout := context.Background()
	/*ctxTimeout, cancel := context.WithTimeout(ctx, m.dbTimeout)
	defer cancel()*/

	var err error

	stmt, err := m.executor.Get().(*sql.Tx).PrepareContext(ctxTimeout, query)
	if err != nil {
		_, _ = m.executor.Get().(*sql.Tx).PrepareContext(ctxTimeout, query)
		panic(err)
		//return fmt.Errorf("error trying to persist product on database: %v", err)
	}
	defer func() {
		err2 := stmt.Close()
		if err2 != nil {
			panic(err2)
		}
	}()

	_, err = stmt.ExecContext(ctxTimeout, args...)
	if err != nil {
		panic(err)
		//return fmt.Errorf("error trying to persist product on database: %v", err)
	}

	return nil
}

func (m ProductRepository) Of(ctx context.Context, id *domain.ProductId) (*domain.Product, error) {
	productSQLStruct := sqlbuilder.NewStruct(new(sqlProduct))

	sb := productSQLStruct.SelectFrom(sqlProductTable)
	sb.Where(sb.Equal("id", id.String()))

	ctxTimeout, cancel := context.WithTimeout(ctx, m.dbTimeout)
	defer cancel()

	// Execute the query.
	sb.Limit(1)
	rawSql, args := sb.Build()

	var rows *sql.Rows
	var err error

	switch m.executor.Get().(type) {
	case *sql.Tx:
		rows, err = m.executor.Get().(*sql.Tx).QueryContext(ctxTimeout, rawSql, args...)
	case *sql.DB:
		rows, err = m.executor.Get().(*sql.DB).QueryContext(ctxTimeout, rawSql, args...)
	}

	if nil != err {
		return &domain.Product{}, fmt.Errorf("error trying to get a query database: %v", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	var product sqlProduct
	for rows.Next() {
		if err = rows.Scan(productSQLStruct.Addr(&product)...); err != nil {
			panic(err)
			//return &domain.Product{}, fmt.Errorf("error trying to get a product on database: %v", err)
		}

		return domain.NewProduct(id.String())
	}

	return &domain.Product{}, fmt.Errorf("%v: %s", domain.ErrProductNotFound, id.String())
}

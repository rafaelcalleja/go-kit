package adapters

import (
	"context"
	"fmt"
	"sync"
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
	connection transaction.SafeQuerier
	dbTimeout  time.Duration
	mux        sync.Mutex
}

func NewMysqlProductRepository(connection transaction.SafeQuerier, dbTimeout time.Duration) *ProductRepository {
	repository := &ProductRepository{
		connection: connection,
		dbTimeout:  dbTimeout,
	}

	return repository
}

func (m *ProductRepository) Save(ctx context.Context, product *domain.Product) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	productSQLStruct := sqlbuilder.NewStruct(new(sqlProduct))
	query, args := productSQLStruct.InsertInto(sqlProductTable, sqlProduct{
		ID: product.ID().String(),
	}).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, m.dbTimeout)
	defer cancel()

	stmt, err := m.connection.Get(ctx).PrepareContext(ctxTimeout, query)
	if err != nil {
		return fmt.Errorf("error trying to persist product on database: %v", err)
	}

	_, err = stmt.ExecContext(ctxTimeout, args...)

	if err != nil {
		return fmt.Errorf("error trying to persist product on database: %v", err)
	}

	return nil
}

func (m *ProductRepository) Of(ctx context.Context, id *domain.ProductId) (*domain.Product, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	productSQLStruct := sqlbuilder.NewStruct(new(sqlProduct))

	sb := productSQLStruct.SelectFrom(sqlProductTable)
	sb.Where(sb.Equal("id", id.String()))

	ctxTimeout, cancel := context.WithTimeout(ctx, m.dbTimeout)
	defer cancel()

	sb.Limit(1)
	rawSql, args := sb.Build()

	rows, err := m.connection.Get(ctx).QueryContext(ctxTimeout, rawSql, args...)

	if nil != err {
		return &domain.Product{}, fmt.Errorf("error trying to get a query database: %v", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	var product sqlProduct
	for rows.Next() {
		if err = rows.Scan(productSQLStruct.Addr(&product)...); err != nil {
			return &domain.Product{}, fmt.Errorf("error trying to get a product on database: %v", err)
		}

		return domain.NewProduct(product.ID)
	}

	return &domain.Product{}, fmt.Errorf("%v: %s", domain.ErrProductNotFound, id.String())
}

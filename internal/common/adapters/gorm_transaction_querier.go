package adapters

import (
	"context"
	"database/sql"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
	"gorm.io/gorm"
)

type GormDB struct {
	transaction.Querier
	*gorm.DB
}

func (g GormDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return g.DB.Statement.ConnPool.ExecContext(ctx, query, args...)
}

func (g GormDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return g.DB.Statement.ConnPool.PrepareContext(ctx, query)
}

func (g GormDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return g.DB.Statement.ConnPool.QueryContext(ctx, query, args...)
}

func (g GormDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return g.DB.Statement.ConnPool.QueryRowContext(ctx, query)
}

func (g GormDB) Rollback() error {
	return g.DB.Rollback().Error
}

func (g GormDB) Commit() error {
	return g.DB.Commit().Error
}

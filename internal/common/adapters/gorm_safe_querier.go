package adapters

import (
	"context"
	"gorm.io/gorm"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type GormSafeQuerier struct {
	transaction.SafeQuerier
}

func (g *GormSafeQuerier) Get(ctx context.Context) transaction.Querier {
	return g.SafeQuerier.Get(ctx)
}

func (g *GormSafeQuerier) GetGorm(ctx context.Context) *gorm.DB {
	return g.SafeQuerier.Get(ctx).(*GormDB).DB
}

func (g *GormSafeQuerier) GetConnPool(ctx context.Context) gorm.ConnPool {
	return g.SafeQuerier.Get(ctx).(*GormDB).DB.Statement.ConnPool
}

func NewGormSafeQuerier(querier transaction.SafeQuerier) *GormSafeQuerier {
	return &GormSafeQuerier{
		SafeQuerier: querier,
	}
}

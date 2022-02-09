package adapters

import (
	"context"
	"sync"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type GormInitializer struct {
	db *GormDB
	mu sync.Mutex
}

func (e *GormInitializer) Begin(_ context.Context) (transaction.Transaction, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	tx := e.db.Begin(nil)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &GormDB{
		DB: tx,
	}, nil
}

func NewGormInitializer(db *GormDB) transaction.Initializer {
	conn := &GormInitializer{
		db: db,
	}

	return conn
}

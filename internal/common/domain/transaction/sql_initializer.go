package transaction

import (
	"context"
	"database/sql"
	"sync"
)

type SqlDBInitializer struct {
	db *sql.DB
	mu sync.Mutex
}

func (e *SqlDBInitializer) Begin(ctx context.Context) (Transaction, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.db.Ping(); err != nil {
		return nil, err
	}

	return e.db.BeginTx(ctx, nil)
}

func NewSqlDBInitializer(db *sql.DB) Initializer {
	conn := &SqlDBInitializer{
		db: db,
	}

	return conn
}

package transaction

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrConnNotFound = errors.New("connection not found")
)

type connectionMap struct {
	*sync.Map
	conn Querier
	mux  sync.Mutex
}

func (t *connectionMap) LoadFromCtxOrBase(ctx context.Context) Querier {
	conn, err := t.LoadFrom(ctx)

	if err != nil {
		return t.BaseConnection()
	}

	return conn
}

func (t *connectionMap) LoadFrom(ctx context.Context) (Querier, error) {
	connectionId, _ := sessionIdFromContext(ctx)

	conn, ok := t.Load(connectionId.String())

	if false == ok {
		return nil, ErrConnNotFound
	}

	return conn, nil
}

func (t *connectionMap) BaseConnection() Querier {
	return t.conn
}

func (t *connectionMap) Load(key string) (value Querier, ok bool) {
	conn, ok := t.Map.Load(key)

	if false == ok {
		return nil, ok
	}

	return conn.(Querier), ok
}

func (t *connectionMap) Len() int {
	t.mux.Lock()
	defer t.mux.Unlock()

	length := 0
	t.Range(func(_, _ interface{}) bool {
		length++
		return true
	})

	return length
}

func newConnectionMap(conn Querier, syncMap *sync.Map) *connectionMap {
	return &connectionMap{
		Map:  syncMap,
		conn: conn,
	}
}

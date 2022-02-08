package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/rafaelcalleja/go-kit/uuid"
)

var (
	ErrSessionIdNotFound = errors.New("session id not found")
)

type SessionId struct {
	id uuid.UUID
}

func NewSessionId(id string) (*SessionId, error) {
	idVO, err := uuid.New().Parse(id)

	if nil != err {
		return &SessionId{}, fmt.Errorf("%w: %s", ErrWrongUuid, id)
	}

	return &SessionId{idVO}, nil
}

func NewRandomSessionId() (*SessionId, error) {
	return NewSessionId(uuid.New().String(uuid.New().Create()))
}

func (s *SessionId) String() string {
	return uuid.New().String(s.id)
}

func (s *SessionId) Equals(other *SessionId) bool {
	return other.String() == s.String()
}

func contextWithNewRandomSessionId(ctx context.Context) (context.Context, error) {
	sessionId, err := NewRandomSessionId()

	if err != nil {
		return ctx, fmt.Errorf("%w: %s", err, ErrUnableToGenerateNewSession.Error())
	}

	return contextWith(ctx, *sessionId)
}

func contextWith(ctx context.Context, id SessionId) (context.Context, error) {
	return context.WithValue(ctx, transactionKey{}, id.String()), nil
}

func sessionIdFromContext(ctx context.Context) (id *SessionId, err error) {
	idFromCtx := ctx.Value(transactionKey{})

	if nil == idFromCtx {
		return &SessionId{}, ErrSessionIdNotFound
	}

	return NewSessionId(idFromCtx.(string))
}

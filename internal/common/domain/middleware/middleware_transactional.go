package middleware

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
)

type Transactional struct {
	session transaction.TransactionalSession
}

func NewMiddlewareTransactional(session transaction.TransactionalSession) Transactional {
	return Transactional{
		session: session,
	}
}

func (t Transactional) Handle(stack StackMiddleware, ctx Context) error {
	return t.session.ExecuteAtomically(GetDefaultContext(ctx).Ctx, func(_ context.Context) error {
		return stack.Next().Handle(stack, ctx)
	})
}

package middleware

import "github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"

type Transactional struct {
	session transaction.TransactionalSession
}

func NewMiddlewareTransactional(session transaction.TransactionalSession) Transactional {
	return Transactional{
		session: session,
	}
}

func (t Transactional) Handle(stack StackMiddleware, ctx Context) error {
	return t.session.ExecuteAtomically(func() error {
		return stack.Next().Handle(stack, ctx)
	})
}

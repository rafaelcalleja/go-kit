package events

import "context"

type HandlerFn func(ctx context.Context, event Event) error

type FuncHandler struct {
	handlerFn HandlerFn
}

func (f FuncHandler) Handle(ctx context.Context, event Event) error {
	return f.handlerFn(ctx, event)
}

func NewFuncHandler(fn HandlerFn) FuncHandler {
	return FuncHandler{
		fn,
	}
}

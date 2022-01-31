package events

import "context"

type HandlerFn func(ctx context.Context, event Event) error
type IsSubscribeToFn func(event Event) bool

type FuncHandler struct {
	handlerFn       HandlerFn
	isSubscribeToFn IsSubscribeToFn
}

func (f FuncHandler) Handle(ctx context.Context, event Event) error {
	return f.handlerFn(ctx, event)
}

func (f FuncHandler) IsSubscribeTo(event Event) bool {
	return f.isSubscribeToFn(event)
}

func NewFuncHandler(fn HandlerFn, fn2 IsSubscribeToFn) *Handler {
	var h Handler = &FuncHandler{
		fn,
		fn2,
	}

	return &h
}

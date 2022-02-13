package commands

import (
	"context"
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "commands/middleware context value " + k.name
}

var (
	pipelineContextKey = &contextKey{"pipeline_context"}
)

type PipelineContext struct {
	Command Command
	Handler Handler
}

func GetPipelineContext(ctx context.Context) PipelineContext {
	return ctx.Value(pipelineContextKey).(PipelineContext)
}

func withPipelineContext(ctx context.Context, handler Handler, cmd Command) context.Context {
	pipelineContext := PipelineContext{
		cmd,
		handler,
	}

	return context.WithValue(ctx, pipelineContextKey, pipelineContext)
}

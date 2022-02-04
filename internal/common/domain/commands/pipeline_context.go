package commands

import (
	"context"
)

const pipelineContextKey string = "pipeline_context"

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

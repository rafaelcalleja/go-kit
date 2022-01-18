package commands

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/middleware"
)

const pipelineContextKey string = "pipeline_context"

type PipelineContext struct {
	Ctx     context.Context
	Command Command
	Handler Handler
}

func GetPipelineContext(context middleware.Context) PipelineContext {
	return context.Get(pipelineContextKey).(PipelineContext)
}

func setPipelineContext(handler Handler, ctx context.Context, cmd Command, middlewareCtx middleware.Context) {
	pipelineContext := PipelineContext{
		ctx,
		cmd,
		handler,
	}

	middlewareCtx.Set(pipelineContextKey, pipelineContext)
}

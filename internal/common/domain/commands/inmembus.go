package commands

import (
	"context"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/middleware"
	"github.com/rafaelcalleja/go-kit/uuid"
)

type busKey struct {
	name string
}

func (c *busKey) String() string {
	return "commands/inmembus context value" + c.name
}

var (
	ctxCommandIdKey = &busKey{"bus_id"}
)

func GetCommandId(ctx context.Context) string {
	id, _ := ctx.Value(ctxCommandIdKey).(string)
	return id
}

func WithCommandId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxCommandIdKey, id)
}

// CommandBus is an in-memory implementation of the commands.Bus.
type CommandBus struct {
	handlers map[Type]Handler
	pipeline Pipeline
}

func NewInMemCommandBusWith(options ...func(*CommandBus) error) (CommandBus, error) {
	var commandBus = new(CommandBus)

	for _, option := range options {
		err := option(commandBus)
		if err != nil {
			return CommandBus{}, err
		}
	}

	if nil == commandBus.handlers {
		commandBus.handlers = make(map[Type]Handler)
	}

	return *commandBus, nil
}

func InMemCommandBusWithPipeline(pipeline Pipeline) func(*CommandBus) error {
	return func(s *CommandBus) error {
		s.pipeline = pipeline
		return nil
	}
}

// NewInMemCommandBus initializes a new instance of CommandBus.
func NewInMemCommandBus() Bus {
	pipeline := NewPipeline()

	commandBus, _ := NewInMemCommandBusWith(
		InMemCommandBusWithPipeline(pipeline),
	)

	return &commandBus
}

// Dispatch implements the commands.Bus interface.
func (b *CommandBus) Dispatch(ctx context.Context, cmd Command) error {
	ctx = WithCommandId(ctx, uuid.New().String(uuid.New().Create()))

	handler, ok := b.handlers[cmd.Type()]
	if !ok {
		return nil
	}

	return b.handle(ctx, handler, cmd)
}

func (b *CommandBus) handle(ctx context.Context, handler Handler, cmd Command) error {
	return b.pipeline.Handle(ctx, handler, cmd)
}

func (b *CommandBus) UseMiddleware(middleware ...middleware.Middleware) {
	b.pipeline.Add(middleware...)
}

func (b *CommandBus) ResetMiddleware() {
	b.pipeline = NewPipeline()
}

// Register implements the commands.Bus interface.
func (b *CommandBus) Register(cmdType Type, handler Handler) {
	b.handlers[cmdType] = handler
}

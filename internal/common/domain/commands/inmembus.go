package commands

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/middleware"
	"sync"
)

type commandBusKey string

func (c commandBusKey) String() string {
	return "command_bus_" + string(c)
}

var (
	ContextDispatchingCommand = commandBusKey("dispatching")
)

// CommandBus is an in-memory implementation of the commands.Bus.
type CommandBus struct {
	handlers map[Type]Handler
	pipeline Pipeline
	mu       sync.Mutex
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
	b.mu.Lock()
	defer b.mu.Unlock()

	handler, ok := b.handlers[cmd.Type()]
	if !ok {
		return nil
	}

	return b.handle(handler, ctx, cmd)
}

func (b *CommandBus) handle(handler Handler, ctx context.Context, cmd Command) error {
	return b.pipeline.Handle(handler, context.WithValue(ctx, ContextDispatchingCommand, cmd), cmd)
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

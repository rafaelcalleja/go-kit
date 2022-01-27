package service

import (
	"context"

	"github.com/rafaelcalleja/go-kit/internal/common/domain/commands"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/events"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/queries"
	"github.com/rafaelcalleja/go-kit/internal/store/application"
	"github.com/rafaelcalleja/go-kit/internal/store/application/command"
	"github.com/rafaelcalleja/go-kit/internal/store/application/event"
	"github.com/rafaelcalleja/go-kit/internal/store/domain"
)

func NewApplication(
	ctx context.Context,
	productRepository domain.ProductRepository,
	commandBus commands.Bus,
	queryBus queries.Bus,
	eventBus events.Bus,
) application.Application {
	creteProductHandler := command.NewCreateProductHandler(productRepository, eventBus)

	commandBus.Register(command.CreateProductCommandType, creteProductHandler)

	eventStore := events.NewInMemEventStore()

	eventBus.Subscribe(
		events.NewStoreEventsOnEventCreated(
			eventStore,
		),
	)

	eventBus.Subscribe(
		event.NewIncreaseStockOnProductCreated(
			domain.NewStockCreateService(),
		),
	)

	return application.Application{
		CommandBus: commandBus,
		QueryBus:   queryBus,
	}
}

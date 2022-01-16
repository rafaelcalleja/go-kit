package main

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/commands"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/events"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/queries"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
	"github.com/rafaelcalleja/go-kit/internal/common/genproto/store"
	"github.com/rafaelcalleja/go-kit/internal/common/server"
	"github.com/rafaelcalleja/go-kit/internal/store/mock"
	"github.com/rafaelcalleja/go-kit/internal/store/ports"
	"github.com/rafaelcalleja/go-kit/internal/store/service"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	commandBus := commands.NewTransactionalCommandBus(
		commands.NewInMemCommandBus(),
		transaction.NewTransactionalSession(
			transaction.NewMockInitializer(),
		),
	)

	queryBus := queries.NewInMemQueryBus()
	eventBus := events.NewInMemoryEventBus()

	application := service.NewApplication(ctx, mock.NewMockProductRepository(), commandBus, queryBus, eventBus)

	server.RunGRPCServer(func(server *grpc.Server) {
		svc := ports.NewGrpcServer(application)
		store.RegisterStoreServiceServer(server, svc)
	})
}

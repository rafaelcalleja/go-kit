package service

import (
	"context"
	"errors"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/commands"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/events"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/queries"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
	"github.com/rafaelcalleja/go-kit/internal/common/genproto/store"
	"github.com/rafaelcalleja/go-kit/internal/common/server"
	"github.com/rafaelcalleja/go-kit/internal/common/tests"
	"github.com/rafaelcalleja/go-kit/internal/store/domain"
	"github.com/rafaelcalleja/go-kit/internal/store/mock"
	"github.com/rafaelcalleja/go-kit/internal/store/ports"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"log"
	"os"
	"testing"
)

var (
	grpcAddr          = "localhost:3000"
	productRepository = mock.NewMockProductRepository()
	commandBus        = commands.NewTransactionalCommandBus(
		commands.NewInMemCommandBus(),
		transaction.NewTransactionalSession(
			transaction.NewMockInitializer(),
		),
	)
	queryBus = queries.NewInMemQueryBus()
	eventBus = events.NewInMemoryEventBus()
)

func TestGrpcClientCreatingProduct(t *testing.T) {
	t.Parallel()

	client := tests.NewStoreGrpcClient(t, grpcAddr)
	productId := "c4546c87-c699-42cb-967a-73a99cd9b7c9"

	called := 0
	productRepository.OfFn = func(id *domain.ProductId) (*domain.Product, error) {
		called++
		return &domain.Product{}, errors.New("product not found")
	}

	productRepository.SaveFn = func(p *domain.Product) error {
		called++
		return nil
	}

	err := client.CreateProduct(context.Background(), productId)
	require.NoError(t, err)
	require.Equal(t, called, 2)
}

func startService() bool {
	app := NewApplication(context.Background(), productRepository, commandBus, queryBus, eventBus)

	go server.RunGRPCServerOnAddr(grpcAddr, func(server *grpc.Server) {
		svc := ports.NewGrpcServer(app)
		store.RegisterStoreServiceServer(server, svc)
	})

	ok := tests.WaitForPort(grpcAddr)
	if !ok {
		log.Println("Timed out waiting for store Grpc to come up")
	}

	return ok
}

func TestMain(m *testing.M) {
	if !startService() {
		log.Println("Timed out waiting for Grpc to come up")
		os.Exit(1)
	}

	os.Exit(m.Run())
}

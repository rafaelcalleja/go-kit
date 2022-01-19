package service

import (
	"context"
	"fmt"
	"github.com/rafaelcalleja/go-kit/internal/common/tests/mysql_tests"
	"log"
	"os"
	"testing"
	"time"

	common_adapters "github.com/rafaelcalleja/go-kit/internal/common/adapters"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/commands"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/events"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/middleware"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/queries"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
	"github.com/rafaelcalleja/go-kit/internal/common/genproto/store"
	"github.com/rafaelcalleja/go-kit/internal/common/server"
	"github.com/rafaelcalleja/go-kit/internal/common/tests"
	"github.com/rafaelcalleja/go-kit/internal/store/adapters"
	"github.com/rafaelcalleja/go-kit/internal/store/domain"
	"github.com/rafaelcalleja/go-kit/internal/store/ports"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var (
	grpcAddr           = "localhost:3000"
	mysqlConnection, _ = mysql_tests.NewMySQLConnection()
	productRepository  = adapters.NewMysqlProductRepository(mysqlConnection.DB, 60*time.Second)
	inMemBus           = commands.NewInMemCommandBus()
	commandBus         = commands.NewTransactionalCommandBus(
		inMemBus,
		transaction.NewTransactionalSession(
			common_adapters.NewTransactionInitializerDb(mysqlConnection),
		),
	)
	queryBus = queries.NewInMemQueryBus()
	eventBus = events.NewInMemoryEventBus()
)

func TestGrpcClientCreatingProduct(t *testing.T) {
	t.Parallel()

	mysql_tests.TruncateTables(mysqlConnection.DB, []string{"products", "stock_products"})

	productId := "c4546c87-c699-42cb-967a-73a99cd9b7c9"

	eventCalledCounter := 0
	eventBus.Subscribe(
		domain.ProductCreatedEventType,
		events.NewFuncHandler(func(ctx context.Context, event events.Event) error {
			eventCalledCounter++
			require.Equal(t, event.AggregateID(), productId)
			return nil
		}),
	)

	inMemBus.(*commands.CommandBus).UseMiddleware(
		middleware.NewMiddlewareTransactional(transaction.NewTransactionalSession(
			transaction.NewMockInitializer(),
		)),
		middleware.NewMiddlewareFunc(func(stack middleware.StackMiddleware, ctx middleware.Context) error {
			pipelineContext := commands.GetPipelineContext(ctx)

			defer fmt.Printf("POST - execute %s\n", pipelineContext.Command.Type())
			fmt.Printf("PRE - execute %s\n", pipelineContext.Command.Type())
			return stack.Next().Handle(stack, ctx)
		}),
	)

	client := tests.NewStoreGrpcClient(t, grpcAddr)

	ctx := context.Background()
	err := client.CreateProduct(ctx, productId)
	require.NoError(t, err)
	require.Equal(t, eventCalledCounter, 1)

	p, _ := domain.NewProductId(productId)

	productFromRepository, err := productRepository.Of(ctx, p)
	require.NoError(t, err)
	require.Equal(t, productFromRepository.ID().String(), p.String())
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

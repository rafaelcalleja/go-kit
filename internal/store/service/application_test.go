package service

import (
	"context"
	"fmt"
	"github.com/rafaelcalleja/go-kit/internal/common/tests/mysql_tests"
	"log"
	"os"
	"sync"
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
	"google.golang.org/grpc/metadata"
)

var (
	grpcAddr           = "localhost:3000"
	mysqlConnection, _ = mysql_tests.NewMySQLConnection()
	executor           = transaction.NewExecutor()
	productRepository  = adapters.NewMysqlProductRepository(executor.Set(mysqlConnection.DB), 60*time.Second)
	inMemBus           = commands.NewInMemCommandBus()
	commandBus         = commands.NewTransactionalCommandBus(
		inMemBus,
		transaction.NewTransactionalSession(
			transaction.NewChainTxInitializer(
				common_adapters.NewTransactionInitializerExecutorSimpleDb(mysqlConnection.DB, executor),
				events.NewMementoTx(eventStoreA),
				events.NewMementoTx(eventStoreB),
			),
		),
	)
	queryBus    = queries.NewInMemQueryBus()
	eventBus    = events.NewInMemoryEventBus()
	eventStoreA = events.NewInMemEventStore()
	eventStoreB = events.NewInMemEventStore()
	lock        = sync.Mutex{}
)

func TestGrpcClientCreatingProduct(t *testing.T) {
	t.Parallel()

	productId := "c4546c87-c699-42cb-967a-73a99cd9b7c9"

	eventCalledCounter := 0
	eventBus.Subscribe(
		events.NewFuncHandler(func(ctx context.Context, event events.Event) error {
			lock.Lock()
			defer lock.Unlock()
			md, _ := metadata.FromIncomingContext(ctx)
			if len(md.Get("testId")) > 0 && md.Get("testId")[0] == "TestGrpcClientCreatingProduct" {
				eventCalledCounter++
			}
			return nil
		}, func(event events.Event) bool {
			return event.Type() == domain.ProductCreatedEventType
		}),
	)

	eventBus.Subscribe(
		events.NewStoreEventsOnEventCreated(
			eventStoreA,
		),
	)

	inMemBus.(*commands.CommandBus).UseMiddleware(
		middleware.NewMiddlewareTransactional(transaction.NewTransactionalSession(
			transaction.NewMockInitializer(),
		)),
		middleware.NewMiddlewareFunc(func(stack middleware.StackMiddleware, ctx middleware.Context) error {
			pipelineContext := commands.GetPipelineContext(ctx)

			defer func() {
				fmt.Printf("POST - execute %s\n", pipelineContext.Command.Type())
			}()

			fmt.Printf("PRE - execute %s\n", pipelineContext.Command.Type())
			return stack.Next().Handle(stack, ctx)
		}),
	)

	client := tests.NewStoreGrpcClient(t, grpcAddr)

	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(
		ctx,
		metadata.Pairs("testId", "TestGrpcClientCreatingProduct"),
	)

	err := client.CreateProduct(ctx, productId)
	require.NoError(t, err)
	require.Equal(t, eventCalledCounter, 1)

	p, _ := domain.NewProductId(productId)

	productFromRepository, err := productRepository.Of(ctx, p)
	require.NoError(t, err)
	require.Equal(t, productFromRepository.ID().String(), p.String())

	evt := eventStoreA.Events()
	require.Equal(t, 1, len(evt))
	require.Equal(t, evt[0].Type(), domain.ProductCreatedEventType)
	require.Equal(t, evt[0].AggregateID(), productId)
}

func TestGrpcClientPanicCreatingProduct(t *testing.T) {
	t.Parallel()

	productId := "ee0dc09e-f3c2-454d-9e25-878d3637a3e4"

	eventBus.Subscribe(
		events.NewStoreEventsOnEventCreated(
			eventStoreB,
		),
	)

	inMemBus.(*commands.CommandBus).UseMiddleware(
		middleware.NewMiddlewareFunc(func(stack middleware.StackMiddleware, ctx middleware.Context) error {
			defer func() {
				lock.Lock()
				defer lock.Unlock()
				ctx2 := commands.GetPipelineContext(ctx).Ctx
				md, _ := metadata.FromIncomingContext(ctx2)

				if len(md.Get("testId")) > 0 && md.Get("testId")[0] == "TestGrpcClientPanicCreatingProduct" {
					panic("fooBar")
				}
			}()
			return stack.Next().Handle(stack, ctx)
		}),
	)

	client := tests.NewStoreGrpcClient(t, grpcAddr)

	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(
		ctx,
		metadata.Pairs("testId", "TestGrpcClientPanicCreatingProduct"),
	)

	err := client.CreateProduct(ctx, productId)
	require.Error(t, err)
	require.Equal(t, "rpc error: code = Internal desc = panic in operation: fooBar", err.Error())

	p, _ := domain.NewProductId(productId)
	_, err = productRepository.Of(ctx, p)
	require.Error(t, err)
	require.Equal(t, "product not found: ee0dc09e-f3c2-454d-9e25-878d3637a3e4", err.Error())

	require.Equal(t, 0, len(eventStoreB.Events()))
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

	mysql_tests.TruncateTables(mysqlConnection.DB, []string{"products", "stock_products"})

	os.Exit(m.Run())
}

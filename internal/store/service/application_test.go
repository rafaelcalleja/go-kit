package service

import (
	"context"
	"fmt"
	"github.com/rafaelcalleja/go-kit/internal/common/tests/mysql_tests"
	"github.com/rafaelcalleja/go-kit/internal/store/application/command"
	"github.com/rafaelcalleja/go-kit/uuid"
	"log"
	"math/rand"
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
	productRepository  = adapters.NewMysqlProductRepository(executor.Set(mysqlConnection), 60*time.Second)
	inMemBus           = commands.NewInMemCommandBus()
	waitChannel        = make(chan bool, 1)
	commandBus         = commands.NewWaiterBus(
		commands.NewTransactionalCommandBus(
			inMemBus,
			transaction.NewTransactionalSession(
				common_adapters.NewTransactionInitializerExecutorSimpleDb(mysqlConnection, executor),
				/*transaction.NewChainTxInitializer(
					common_adapters.NewTransactionInitializerExecutorChanDb(mysqlConnection.DB, executor),
					events.NewMementoTx(eventStoreA),
					events.NewMementoTx(eventStoreB),
				),*/
			),
		),
		waitChannel,
	)

	queryBus    = queries.NewInMemQueryBus()
	eventBus    = events.NewInMemoryEventBus()
	eventStoreA = events.NewInMemEventStore()
	eventStoreB = events.NewInMemEventStore()
)

func aTestTryDead(t *testing.T) {
	t.Parallel()
	wg := sync.WaitGroup{}

	u := uuid.New()
	client := tests.NewStoreGrpcClient(t, grpcAddr)

	min := 1
	max := 3
	fmt.Println(rand.Intn(max-min) + min)

	inMemBus.(*commands.CommandBus).UseMiddleware(
		middleware.NewMiddlewareFunc(func(stack middleware.StackMiddleware, ctx middleware.Context) error {
			defer func() {
				ctx2 := commands.GetPipelineContext(ctx).Ctx
				md, _ := metadata.FromIncomingContext(ctx2)
				if len(md.Get("testId")) > 0 && md.Get("testId")[0] == "TestTryDead" && rand.Intn(max-min)+min == 1 {
					panic("fooBar")
				}
			}()
			return stack.Next().Handle(stack, ctx)
		}),
	)

	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(
		ctx,
		metadata.Pairs("testId", "TestTryDead"),
	)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			conn, _ := mysql_tests.NewMySQLConnection()
			executor.Set(conn)
			productId := u.String(u.Create())
			err := client.CreateProduct(ctx, productId)
			if err != nil {
				require.Equal(t, "rpc error: code = Internal desc = panic in operation: fooBar", err.Error())
			}

			wg.Done()
		}()
	}

	wg.Wait()
}

func TestGrpcClientCreatingProduct2(t *testing.T) {
	t.Parallel()

	productId := "c4546c87-c699-42cb-967a-73a99cd9b7c8"

	eventCalledCounter := 0
	eventBus.Subscribe(
		events.NewFuncHandler(func(ctx context.Context, event events.Event) error {
			md, _ := metadata.FromIncomingContext(ctx)
			if len(md.Get("testId")) > 0 && md.Get("testId")[0] == "TestGrpcClientCreatingProduct2" {
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
		/*middleware.NewMiddlewareTransactional(transaction.NewTransactionalSession(
			transaction.NewMockInitializer(),
		)),*/
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
		metadata.Pairs("testId", "TestGrpcClientCreatingProduct2"),
	)

	err := client.CreateProduct(ctx, productId)
	require.NoError(t, err)
	require.GreaterOrEqual(t, eventCalledCounter, 1)

	p, _ := domain.NewProductId(productId)

	productFromRepository, err := productRepository.Of(ctx, p)
	require.NoError(t, err)
	require.Equal(t, productFromRepository.ID().String(), p.String())

	evt := eventStoreA.Events()
	require.GreaterOrEqual(t, len(evt), 1)
	/*require.Equal(t, evt[0].Type(), domain.ProductCreatedEventType)
	require.Equal(t, evt[0].AggregateID(), productId)*/
}

func TestGrpcClientCreatingProduct(t *testing.T) {
	t.Parallel()

	productId := "c4546c87-c699-42cb-967a-73a99cd9b7c9"

	eventCalledCounter := 0
	eventBus.Subscribe(
		events.NewFuncHandler(func(ctx context.Context, event events.Event) error {
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
		/*middleware.NewMiddlewareTransactional(transaction.NewTransactionalSession(
			transaction.NewMockInitializer(),
		)),*/
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
	require.GreaterOrEqual(t, eventCalledCounter, 1)

	p, _ := domain.NewProductId(productId)

	productFromRepository, err := productRepository.Of(ctx, p)
	require.NoError(t, err)
	require.Equal(t, productFromRepository.ID().String(), p.String())

	evt := eventStoreA.Events()
	require.GreaterOrEqual(t, len(evt), 1)
	/*require.Equal(t, evt[0].Type(), domain.ProductCreatedEventType)
	require.Equal(t, evt[0].AggregateID(), productId)*/
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

	//	require.Equal(t, 0, len(eventStoreB.Events()))
	err := client.CreateProduct(ctx, productId)
	require.Error(t, err)
	require.Equal(t, "rpc error: code = Internal desc = panic in operation: fooBar", err.Error())

	p, _ := domain.NewProductId(productId)
	_, err = productRepository.Of(ctx, p)
	require.Error(t, err)
	require.Equal(t, "product not found: ee0dc09e-f3c2-454d-9e25-878d3637a3e4", err.Error())

	/*for _, e := range eventStoreB.Events() {
		require.NotEqual(t, "ee0dc09e-f3c2-454d-9e25-878d3637a3e4", e.AggregateID())
	}*/
	//require.Equal(t, 0, len(eventStoreB.Events()))
}
func TestGrpcClientPanicCreatingProduct2(t *testing.T) {
	t.Parallel()

	productId := "ee0dc09e-f3c2-454d-9e25-878d3637a3e5"

	eventBus.Subscribe(
		events.NewStoreEventsOnEventCreated(
			eventStoreB,
		),
	)

	inMemBus.(*commands.CommandBus).UseMiddleware(
		middleware.NewMiddlewareFunc(func(stack middleware.StackMiddleware, ctx middleware.Context) error {
			defer func() {
				ctx2 := commands.GetPipelineContext(ctx).Ctx
				md, _ := metadata.FromIncomingContext(ctx2)

				if len(md.Get("testId")) > 0 && md.Get("testId")[0] == "TestGrpcClientPanicCreatingProduct2" {
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
		metadata.Pairs("testId", "TestGrpcClientPanicCreatingProduct2"),
	)

	require.Equal(t, 0, len(eventStoreB.Events()))
	err := client.CreateProduct(ctx, productId)
	require.Error(t, err)
	require.Equal(t, "rpc error: code = Internal desc = panic in operation: fooBar", err.Error())

	p, _ := domain.NewProductId(productId)
	_, err = productRepository.Of(ctx, p)
	require.Error(t, err)
	require.Equal(t, "product not found: ee0dc09e-f3c2-454d-9e25-878d3637a3e5", err.Error())

	/*for _, e := range eventStoreB.Events() {
		require.NotEqual(t, "ee0dc09e-f3c2-454d-9e25-878d3637a3e4", e.AggregateID())
	}*/
	//require.Equal(t, 0, len(eventStoreB.Events()))
}

func TestGrpcClientPanicCreatingProduct3(t *testing.T) {
	t.Parallel()

	productId := "ee0dc09e-f3c2-454d-9e25-878d3637a3e6"

	eventBus.Subscribe(
		events.NewStoreEventsOnEventCreated(
			eventStoreB,
		),
	)

	inMemBus.(*commands.CommandBus).UseMiddleware(
		middleware.NewMiddlewareFunc(func(stack middleware.StackMiddleware, ctx middleware.Context) error {
			defer func() {
				ctx2 := commands.GetPipelineContext(ctx).Ctx
				md, _ := metadata.FromIncomingContext(ctx2)

				if len(md.Get("testId")) > 0 && md.Get("testId")[0] == "TestGrpcClientPanicCreatingProduct3" {
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
		metadata.Pairs("testId", "TestGrpcClientPanicCreatingProduct3"),
	)

	require.Equal(t, 0, len(eventStoreB.Events()))
	err := client.CreateProduct(ctx, productId)
	require.Error(t, err)
	require.Equal(t, "rpc error: code = Internal desc = panic in operation: fooBar", err.Error())

	p, _ := domain.NewProductId(productId)
	_, err = productRepository.Of(ctx, p)
	require.Error(t, err)
	require.Equal(t, "product not found: ee0dc09e-f3c2-454d-9e25-878d3637a3e6", err.Error())

	/*for _, e := range eventStoreB.Events() {
		require.NotEqual(t, "ee0dc09e-f3c2-454d-9e25-878d3637a3e4", e.AggregateID())
	}*/
	//require.Equal(t, 0, len(eventStoreB.Events()))
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

	waitChannel <- true

	return ok
}

func TestMain(m *testing.M) {
	if !startService() {
		log.Println("Timed out waiting for Grpc to come up")
		os.Exit(1)
	}

	mysql_tests.TruncateTables(mysqlConnection, []string{"products", "stock_products"})

	os.Exit(m.Run())
}

func TestChannels(t *testing.T) {
	wait := waitChannel
	conn, _ := mysql_tests.NewMySQLConnection()
	_ = common_adapters.NewTransactionInitializerExecutorChanDb(conn, wait, executor)
	ctx := context.Background()
	u := uuid.New()
	client := tests.NewStoreGrpcClient(t, grpcAddr)

	for x := 0; x < 10; x++ {
		t.Run("parallel test", func(t *testing.T) {
			t.Parallel()

			size := 20
			wg := sync.WaitGroup{}
			wg.Add(size)

			for i := 0; i < size; i++ {
				go func() {
					/*_, _ = initializer.Begin()
					stmt, _ := executor.Get().(*sql.Tx).Prepare("SELECT 1")
					_, _ = stmt.Exec(nil)
					err := initializer.Rollback()
					require.NoError(t, err)*/

					productId := u.String(u.Create())
					_ = client.CreateProduct(ctx, productId)
					_ = client.CreateProduct(ctx, productId)
					err := commandBus.Dispatch(ctx, command.NewCreateProductCommand(productId))
					require.Error(t, err)
					wg.Done()
				}()
			}

			wg.Wait()
		})
	}
}

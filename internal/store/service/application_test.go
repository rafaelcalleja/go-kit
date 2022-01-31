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
	"github.com/rafaelcalleja/go-kit/internal/common/domain/queries"
	common_sync "github.com/rafaelcalleja/go-kit/internal/common/domain/sync"
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
	executor           = transaction.NewExecutor()
	productRepository  = adapters.NewMysqlProductRepository(executor, 60*time.Second, &lock, locker)
	inMemBus           = commands.NewInMemCommandBus()
	waitChannel        = make(chan bool, 1)
	commandBus         = commands.NewWaiterBus(
		commands.NewTransactionalCommandBus(
			inMemBus,
			transaction.NewTransactionalSession(
				transaction.NewChainTxInitializer(
					common_adapters.NewTransactionInitializerExecutorSimpleDb(mysqlConnection, executor),
					events.NewMementoTx(eventStore),
				),
			),
		),
		locker,
	)

	queryBus   = queries.NewInMemQueryBus()
	eventBus   = events.NewInMemoryEventBus()
	eventStore = events.NewInMemEventStore()
	lock       = sync.RWMutex{}
	cond       = sync.NewCond(&lock)
	locker     = common_sync.NewChanSync()
)

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
	ctx := context.Background()
	u := uuid.New()
	client := tests.NewStoreGrpcClient(t, grpcAddr)

	handler := events.NewStoreEventsOnEventCreated(
		eventStore,
	)

	eventBus.Subscribe(
		handler,
	)

	size := 20
	var loops int

	t.Run("group test", func(t *testing.T) {
		for x := 0; x < 10; x++ {
			t.Run("parallel test", func(t *testing.T) {
				loops += 1
				t.Parallel()

				wg := &sync.WaitGroup{}
				wg.Add(size)

				for i := 0; i < size; i++ {
					go func() {
						productId := u.String(u.Create())
						_ = client.CreateProduct(ctx, productId)

						err := client.CreateProduct(ctx, productId)
						require.Error(t, err)
						require.Equal(t, fmt.Sprintf("rpc error: code = Internal desc = product exists: %s", productId), err.Error())

						err = commandBus.Dispatch(ctx, command.NewCreateProductCommand(productId))
						require.Error(t, err)
						require.Equal(t, fmt.Sprintf("product exists: %s", productId), err.Error())

						id, _ := domain.NewProductId(productId)
						_, err = productRepository.Of(ctx, id)
						require.NoError(t, err)

						wg.Done()
					}()
				}

				wg.Wait()

			})
		}
	})

	require.Equal(t, size*loops, len(eventStore.Events()))

	exceptionHandler := events.NewFuncHandler(func(ctx context.Context, event events.Event) error {
		panic(fmt.Sprintf("foo %s", event.AggregateID()))
	}, func(event events.Event) bool {
		return event.Type() == domain.ProductCreatedEventType
	})

	eventBus.Subscribe(
		exceptionHandler,
	)

	t.Run("group exception test", func(t *testing.T) {
		for x := 0; x < 10; x++ {
			t.Run("parallel exception test", func(t *testing.T) {
				t.Parallel()

				wg := &sync.WaitGroup{}
				wg.Add(size)

				for i := 0; i < size; i++ {
					go func() {
						productId := u.String(u.Create())
						err := client.CreateProduct(ctx, productId)
						require.Error(t, err)
						require.Equal(t, fmt.Sprintf("rpc error: code = Internal desc = panic in operation: foo %s", productId), err.Error())

						err = commandBus.Dispatch(ctx, command.NewCreateProductCommand(productId))
						require.Error(t, err)
						require.Equal(t, fmt.Sprintf("panic in operation: foo %s", productId), err.Error())

						time.Sleep(time.Duration(rand.Intn(50000)) * time.Microsecond)
						id, _ := domain.NewProductId(productId)
						_, err = productRepository.Of(ctx, id)
						require.Error(t, err)
						require.Equal(t, fmt.Sprintf("product not found: %s", productId), err.Error())

						wg.Done()
					}()
				}

				wg.Wait()

			})
		}
	})

	eventBus.Unsubscribe(
		exceptionHandler,
	)

	require.Equal(t, size*loops, len(eventStore.Events()))
}

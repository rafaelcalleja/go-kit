package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rafaelcalleja/go-kit/internal/store/adapters"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/commands"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/events"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/queries"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/transaction"
	"github.com/rafaelcalleja/go-kit/internal/common/genproto/store"
	"github.com/rafaelcalleja/go-kit/internal/common/server"
	"github.com/rafaelcalleja/go-kit/internal/store/ports"
	"github.com/rafaelcalleja/go-kit/internal/store/service"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	var cfg mysqlCfg
	err := envconfig.Process("MYSQL", &cfg)
	if err != nil {
		panic(err)
	}

	mysqlURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.User, cfg.Password, cfg.Addr, cfg.Port, cfg.Database)
	mysqlConnection, err := sql.Open("mysql", mysqlURI)

	if err != nil {
		panic(err)
	}

	txHandler := transaction.NewTxHandler(mysqlConnection)
	commandBus := commands.NewTransactionalCommandBus(
		commands.NewInMemCommandBus(),
		transaction.NewTransactionalSession(
			transaction.NewTxHandlerInitializer(txHandler, transaction.NewSqlDBInitializer(mysqlConnection)),
		),
	)

	queryBus := queries.NewInMemQueryBus()
	eventBus := events.NewInMemoryEventBus()
	productRepository := adapters.NewMysqlProductRepository(txHandler.SafeQuerier(), cfg.Timeout)

	application := service.NewApplication(ctx, productRepository, commandBus, queryBus, eventBus)

	go server.RunGRPCServer(func(server *grpc.Server) {
		svc := ports.NewGrpcServer(application)
		store.RegisterStoreServiceServer(server, svc)
	})

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(
			ports.NewHttpServer(application),
			router,
		)
	})
}

type mysqlCfg struct {
	User     string
	Password string
	Addr     string        `default:"localhost"`
	Port     uint          `default:"3306"`
	Database string        `default:"db"`
	Timeout  time.Duration `default:"5s"`
}

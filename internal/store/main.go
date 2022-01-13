package main

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/common/genproto/store"
	"github.com/rafaelcalleja/go-kit/internal/common/server"
	"github.com/rafaelcalleja/go-kit/internal/store/ports"
	"github.com/rafaelcalleja/go-kit/internal/store/service"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	application := service.NewApplication(ctx)

	server.RunGRPCServer(func(server *grpc.Server) {
		svc := ports.NewGrpcServer(application)
		store.RegisterStoreServiceServer(server, svc)
	})
}

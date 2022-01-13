package ports

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rafaelcalleja/go-kit/internal/common/genproto/store"
	"github.com/rafaelcalleja/go-kit/internal/store/application"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	app application.Application
}

func NewGrpcServer(application application.Application) GrpcServer {
	return GrpcServer{app: application}
}

func (g GrpcServer) CreateProduct(ctx context.Context, request *store.CreateProductRequest) (*empty.Empty, error) {
	productId := request.ProductId

	if err := g.app.Commands.CreateProduct.Handle(productId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

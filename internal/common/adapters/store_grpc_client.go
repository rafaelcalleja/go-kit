package adapters

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/common/genproto/store"
)

type StoreGrpc struct {
	client store.StoreServiceClient
}

func NewStoreGrpc(client store.StoreServiceClient) StoreGrpc {
	return StoreGrpc{client: client}
}

func (s StoreGrpc) CreateProduct(ctx context.Context, productId string) error {
	_, err := s.client.CreateProduct(ctx, &store.CreateProductRequest{
		ProductId: productId,
	})

	return err
}

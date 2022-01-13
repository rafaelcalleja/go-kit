package service

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/store/application"
	"github.com/rafaelcalleja/go-kit/internal/store/application/command"
	"github.com/rafaelcalleja/go-kit/internal/store/domain"
)

func NewApplication(ctx context.Context, productRepository domain.ProductRepository) application.Application {
	return application.Application{
		Commands: application.Commands{
			CreateProduct: command.NewCreateProductHandler(productRepository),
		},
	}
}

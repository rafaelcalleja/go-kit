package service

import (
	"context"
	"github.com/rafaelcalleja/go-kit/internal/store/application"
	"github.com/rafaelcalleja/go-kit/internal/store/application/command"
	"github.com/rafaelcalleja/go-kit/internal/store/mock"
)

func NewApplication(ctx context.Context) application.Application {
	return application.Application{
		Commands: application.Commands{
			CreateProduct: command.NewCreateProductHandler(mock.NewMockProductRepository()),
		},
	}
}

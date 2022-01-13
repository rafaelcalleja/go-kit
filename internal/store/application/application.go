package application

import "github.com/rafaelcalleja/go-kit/internal/store/application/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateProduct command.CreateProductHandler
}

type Queries struct {
}

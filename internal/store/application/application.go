package application

import (
	"github.com/rafaelcalleja/go-kit/internal/common/domain/commands"
	"github.com/rafaelcalleja/go-kit/internal/common/domain/queries"
)

type Application struct {
	CommandBus commands.Bus
	QueryBus   queries.Bus
}

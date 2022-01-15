package command

import (
	"github.com/rafaelcalleja/go-kit/internal/common/domain/commands"
)

const CreateProductCommandType commands.Type = "create.product.command"

// CreateProductCommand is the command dispatched to create a new product.
type CreateProductCommand struct {
	id string
}

// NewCreateProductCommand creates a new CourseCommand.
func NewCreateProductCommand(id string) CreateProductCommand {
	return CreateProductCommand{
		id: id,
	}
}

func (c CreateProductCommand) Type() commands.Type {
	return CreateProductCommandType
}

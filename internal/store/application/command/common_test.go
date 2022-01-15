package command

import (
	"github.com/rafaelcalleja/go-kit/internal/common/domain/commands"
)

type wrongCommand struct{}

func newWrongCommand() wrongCommand {
	return wrongCommand{}
}

func (command wrongCommand) Type() commands.Type {
	return "wrong.command.type"
}

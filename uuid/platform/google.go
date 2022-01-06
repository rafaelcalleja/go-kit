package platform

import (
	googleuuid "github.com/google/uuid"
	"github.com/rafaelcalleja/go-kit/uuid"
)

type GoogleUuid struct {
}

func (g *GoogleUuid) Parse(s string) (uuid.UUID, error) {
	guuid, err := googleuuid.Parse(s)

	return uuid.UUID(guuid), err
}

func (g *GoogleUuid) String(uuid uuid.UUID) string {
	guuid := googleuuid.UUID(uuid)

	return guuid.String()
}

func NewGoogle() *GoogleUuid {
	return new(GoogleUuid)
}

package uuid

import (
	googleuuid "github.com/google/uuid"
)

type googleUuid struct {
}

func (g *googleUuid) Parse(s string) (UUID, error) {
	guuid, err := googleuuid.Parse(s)

	return UUID(guuid), err
}

func (g *googleUuid) String(uuid UUID) string {
	guuid := googleuuid.UUID(uuid)

	return guuid.String()
}

func newGoogle() *googleUuid {
	return new(googleUuid)
}

package uuid

import "github.com/rafaelcalleja/go-kit/uuid/platform"

func New() Uuid {
	return platform.NewGoogle()
}

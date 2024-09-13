package writer

import (
	"errors"

	"github.com/fireland15/rpc-gen/internal/compiler"
)

var ErrUndefinedClient = errors.New("no config for client")

type ClientWriter interface {
	Write(service *compiler.Service) error
}

type ServerWriter interface {
	Write(service *compiler.Service) error
}

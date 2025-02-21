package generators

import (
	"errors"
	"log"

	"github.com/fireland15/rpc-gen/internal/config"
	"github.com/fireland15/rpc-gen/internal/generators/typescript"
	"github.com/fireland15/rpc-gen/internal/protocol"
)

var ErrUndefinedClient = errors.New("no config for client")

type CodeGenerator interface {
	Generate(protocol *protocol.Protocol) error
}

type rootGenerator struct {
	inner []CodeGenerator
}

func GeneratorFromConfig(config *config.RpcGenConfig) (CodeGenerator, error) {
	generator := new(rootGenerator)
	if config == nil {
		panic("config is nil")
	}

	for client, clientConfig := range config.Clients {
		log.Printf("Configuring %s client code generator.\n", client)
		if client == "typescript" {
			gen, err := typescript.NewTypescriptGenerator(clientConfig)
			if err != nil {
				return nil, err
			}
			generator.inner = append(generator.inner, gen)
		}
	}

	// for server, serverConfig := range config.Servers {
	// 	log.Printf("Configuring %s server code generator.\n", server)
	// 	if server == "go-echo" {
	// 		gen, err := NewGoEchoServerGenerator(serverConfig)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		generator.inner = append(generator.inner, gen)
	// 	}
	// }

	return generator, nil
}

func (g *rootGenerator) Generate(protocol *protocol.Protocol) error {
	for _, generator := range g.inner {
		err := generator.Generate(protocol)
		if err != nil {
			return err
		}
	}
	return nil
}

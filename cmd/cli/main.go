package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fireland15/rpc-gen/internal/compiler"
	"github.com/fireland15/rpc-gen/internal/config"
	"github.com/fireland15/rpc-gen/internal/generators"
)

func main() {
	configPath := flag.String("c", "config.json", "path to config file")
	config, err := config.ReadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	definitionFile, err := os.Open(config.RpcDefinitionFile)
	if err != nil {
		panic(err)
	}

	c := compiler.Compiler{}
	service, err := c.Compile(definitionFile)
	if err != nil {
		panic(err)
	}

	generator, err := generators.GeneratorFromConfig(config)
	if err != nil {
		err = fmt.Errorf("problem creating generator: %w", err)
		panic(err)
	}

	err = generator.Generate(service)
	if err != nil {
		err = fmt.Errorf("problem generating code for service: %w", err)
		panic(err)
	}

	fmt.Println("Generation complete")
}

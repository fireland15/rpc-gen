package main

import (
	"flag"
	"fmt"

	"github.com/fireland15/rpc-gen/internal/compiler"
	"github.com/fireland15/rpc-gen/internal/config"
)

func main() {
	configPath := flag.String("c", "config.json", "path to config file")

	flag.Parse()

	config, err := config.ReadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	err = compiler.Compile(config.RpcDefinitionFile, config)
	if err != nil {
		panic(err)
	}

	fmt.Println("Generation complete")
}

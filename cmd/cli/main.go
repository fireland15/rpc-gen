package main

import (
	"flag"
	"os"

	"github.com/fireland15/rpc-gen/internal/compiler"
	"github.com/fireland15/rpc-gen/internal/config"
	"github.com/fireland15/rpc-gen/internal/writer"
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

	f, err := os.OpenFile("service_client.ts", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	x := config.Clients["typescript"]
	w := writer.NewTypescriptClientWriter(&x, f)
	err = w.Write(service)
	if err != nil {
		panic(err)
	}

	gof, err := os.OpenFile("server.go", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	sw, err := writer.NewGoServerWriter(config, gof)
	if err != nil {
		panic(err)
	}
	err = sw.Write(service)
	if err != nil {
		panic(err)
	}
}

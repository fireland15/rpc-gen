package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/fireland15/rpc-gen/internal/compiler"
)

func main() {
	configPath := flag.String("c", "config.json", "path to config file")
	flag.Parse()

	cfgBytes, err := os.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	var cfg compiler.Config
	if err := json.Unmarshal(cfgBytes, &cfg); err != nil {
		panic(err)
	}

	if err := compiler.Compile(cfg); err != nil {
		panic(err)
	}

	fmt.Println("Generation complete")
}

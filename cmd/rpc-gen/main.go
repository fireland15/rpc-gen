package main

import (
	"fmt"

	"github.com/fireland15/rpc-gen/internal/compiler"
	"github.com/fireland15/rpc-gen/internal/generator"
)

func main() {
	//configPath := flag.String("c", "config.json", "path to config file")
	//
	//flag.Parse()

	cfg := compiler.Config{
		Generators: map[string]generator.LanguageConfig{
			"typescript": {
				Scalars: map[string]string{
					"String":   "string",
					"DateOnly": "string",
					"Int":      "number",
					"Uuid":     "string",
				},
				Passes: []generator.PassConfig{
					{
						Template: "ts-client",
						Output:   "out/ts-client/client.ts",
					},
					{
						Template: "ts-client",
						Output:   "out/ts-client2/client.ts",
					},
				},
			},
			"cs": {
				Scalars: map[string]string{
					"String":   "string",
					"DateOnly": "DateOnly",
					"Int":      "int",
					"Uuid":     "Guid",
				},
				Passes: []generator.PassConfig{
					{
						Template: "model",
						Output:   "{{.Name}}.cs",
						Dir:      "out/cs/server/models",
						Items:    []string{"composite"},
					},
					{
						Template: "enum",
						Output:   "{{.Name}}.cs",
						Dir:      "out/cs/server/models/enums",
						Items:    []string{"enum"},
					},
					{
						Template: "interface",
						Output:   "I{{toCase .Name \"upperCamel\"}}Handler.cs",
						Dir:      "out/cs/server/abstraction",
						Items:    []string{"method"},
					},
					{
						Template: "minimal_api",
						Output:   "out/cs/server/web/Handlers.cs",
						Items:    []string{"method"},
					},
				},
			},
		},
	}
	err := compiler.Compile("journal.rpc", cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println("Generation complete")
}

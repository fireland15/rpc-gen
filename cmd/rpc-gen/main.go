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
			//"typescript": {
			//	Scalars: map[string]string{
			//		"String":   "string",
			//		"DateOnly": "string",
			//		"Int":      "number",
			//		"Uuid":     "string",
			//	},
			//	Passes: []generator.PassConfig{
			//		{
			//			Template: "ts-client",
			//			Output:   "out/ts-client/client.ts",
			//		},
			//		{
			//			Template: "ts-client",
			//			Output:   "out/ts-client2/client.ts",
			//		},
			//	},
			//},
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
						Dir:      "examples/RpcGen.Examples.Api/RpcGen.Examples/Models",
						Items:    []string{"composite"},
						Inject: map[string]string{
							"namespace": "RpcGen.Examples.Models",
						},
					},
					{
						Template: "methodParameter",
						Output:   "{{.Name}}Parameters.cs",
						Dir:      "examples/RpcGen.Examples.Api/RpcGen.Examples/Models",
						Items:    []string{"hasArgs"},
						Inject: map[string]string{
							"namespace": "RpcGen.Examples.Models",
						},
					},
					{
						Template: "enum",
						Output:   "{{.Name}}.cs",
						Dir:      "examples/RpcGen.Examples.Api/RpcGen.Examples/Models/Enums",
						Items:    []string{"enum"},
						Inject: map[string]string{
							"namespace": "RpcGen.Examples.Models.Enums",
						},
					},
					{
						Template: "interface",
						Output:   "I{{toCase .Name \"upperCamel\"}}Handler.cs",
						Dir:      "examples/RpcGen.Examples.Api/RpcGen.Examples/Endpoints/Abstractions",
						Items:    []string{"method"},
						Inject: map[string]string{
							"namespace": "RpcGen.Examples.Endpoints.Abstractions",
							"usings": `using RpcGen.Examples.Models;
`,
						},
					},
					{
						Template: "minimal_api",
						Output:   "examples/RpcGen.Examples.Api/RpcGen.Examples/Endpoints/Handlers.cs",
						Items:    []string{"method"},
						Inject: map[string]string{
							"namespace": "RpcGen.Examples.Endpoints",
							"usings": `using Microsoft.AspNetCore.Mvc;
using RpcGen.Examples.Endpoints.Abstractions;
using RpcGen.Examples.Models;`,
						},
					},
				},
				Inject: map[string]string{
					"namespace": "RpcGen.Examples",
					"usings": `using System.Text.Json.Serialization;
`,
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

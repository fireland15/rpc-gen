package typescript

import "github.com/fireland15/rpc-gen/internal/model"

type typescriptCodegen struct {
	Imports []importStatement
	Types   model.Model
	Methods model.Method
}

type importStatement struct {
	Symbols []string
	Module  string
}

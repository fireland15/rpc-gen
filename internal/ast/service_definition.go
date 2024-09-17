package ast

import (
	"github.com/fireland15/rpc-gen/internal/lexing"
)

type ServiceDefinition struct {
	MethodDefinitions []MethodDefinition
}

type MethodDefinition struct {
	Name Identifier
}

type ParenthesisGroup[T any] struct {
	Open  lexing.Token
	Close lexing.Token
	Inner T
}

type Identifier struct {
	Name string
	Span lexing.Span
}

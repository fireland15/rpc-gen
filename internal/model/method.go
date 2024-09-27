package model

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

type Method struct {
	Name          string
	Parameters    []MethodParameter
	ReturnType    *Type
	ParameterType Type
}

type MethodParameter struct {
	Name string
	Type Type
}

func (m Method) Path() string {
	return fmt.Sprintf("/%s", strcase.ToSnake(m.Name))
}

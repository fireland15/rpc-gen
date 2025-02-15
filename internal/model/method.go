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
	service       *ProtocolDefinition
}

type MethodParameter struct {
	Name string
	Type Type
}

func (m Method) Path() string {
	return fmt.Sprintf("/%s", strcase.ToSnake(m.Name))
}

func (m Method) HasUpload() bool {
	for _, p := range m.Parameters {
		if p.Type.IsUploadType() {
			return true
		}

		variant, model := m.service.GetTypeDefinition(p.Type.BaseName())
		if variant == UserDefined {
			for _, f := range model.Fields {
				if f.Type.IsUploadType() {
					return true
				}
			}
		}
	}

	return false
}

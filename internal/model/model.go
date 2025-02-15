package model

import "github.com/iancoleman/strcase"

type Model struct {
	Name              string
	Fields            []Field
	serviceDefinition *ProtocolDefinition
}

type Field struct {
	Name string
	Type Type
}

func (f Field) SerializedName() string {
	return strcase.ToLowerCamel(f.Name)
}

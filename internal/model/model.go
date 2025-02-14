package model

import "github.com/iancoleman/strcase"

type Model struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type Type
}

func (f Field) SerializedName() string {
	return strcase.ToLowerCamel(f.Name)
}

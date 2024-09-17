package model

type Model struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type Type
}

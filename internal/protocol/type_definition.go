package protocol

import "errors"

var ErrDuplicateFieldDefinition = errors.New("duplicate field definition")

type TypeDefinitionVariant string
type SerializedType string

const (
	Scalar   TypeDefinitionVariant = "SCALAR"
	Compound TypeDefinitionVariant = "COMPOUND"
)

const (
	SerializedTypeNumber SerializedType = "number"
	SerializedTypeString SerializedType = "string"
	SerializedTypeBool   SerializedType = "bool"
	SerializedTypeObject SerializedType = "object"
	SerializedTypeArray  SerializedType = "array"
	SerializedTypeFile   SerializedType = "File"
)

type TypeDefinition struct {
	Variant        TypeDefinitionVariant
	Name           string
	Fields         []*FieldDefinition
	SerializedType *SerializedType
}

type FieldDefinition struct {
	Name string
	Type *TypeRef2
}

func NewObjectDefinition(name string) *TypeDefinition {
	td := new(TypeDefinition)
	td.Name = name
	td.Fields = make([]*FieldDefinition, 0)
	td.Variant = Compound
	return td
}

func NewScalarDefinition(name string, serializedType SerializedType) *TypeDefinition {
	td := new(TypeDefinition)
	td.Name = name
	td.Fields = make([]*FieldDefinition, 0)
	td.Variant = Scalar
	td.SerializedType = &serializedType
	return td
}

func (td *TypeDefinition) AddField(field *FieldDefinition) error {
	if td.Variant != Compound {
		panic("cannot add fields to scalar type")
	}

	for idx := range td.Fields {
		if td.Fields[idx].Name == field.Name {
			return ErrDuplicateFieldDefinition
		}
	}

	td.Fields = append(td.Fields, field)

	return nil
}

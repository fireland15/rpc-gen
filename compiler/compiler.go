package compiler

import (
	"errors"
	"fmt"
	"io"

	"github.com/fireland15/rpc-gen/parser"
)

type TypeVariant int

const (
	TypeVariantScalar TypeVariant = iota
	TypeVariantObject
	TypeVariantReference
)

type Type struct {
	Name    string
	Variant TypeVariant
	Fields  map[string]Field
}

type Field struct {
	Type     Type
	Optional bool
}

type Rpc struct {
	Name         string
	RequestType  *Type
	ResponseType *Type
}

type Service struct {
	rpc   map[string]Rpc
	types map[string]Type
}

type Compiler struct {
	Errors []error
}

func (c *Compiler) Compile(input io.Reader) (*Service, error) {
	p, err := parser.NewParser(input)
	if err != nil {
		err = fmt.Errorf("failed to init parser: %w", err)
		return nil, err
	}

	def, err := p.Parse()
	if err != nil {
		err = fmt.Errorf("problem parsing input: %w", err)
		return nil, err
	}

	service := new(Service)

}

var ErrMultipleSymbolDefinition = errors.New("symbol defined multiple times")

func (c *Compiler) getDefinedTypes(sd parser.ServiceDefinition) map[string]Type {
	typeMap := make(map[string]Type)

	c.addBuiltInTypes(&typeMap)
	c.addCustomTypes(&sd, &typeMap)

	return typeMap
}

func (c *Compiler) addCustomTypes(sd *parser.ServiceDefinition, typeMap *map[string]Type) {
	for _, model := range sd.Models {
		_, found := (*typeMap)[model.Name]
		if found {
			err := fmt.Errorf("duplicate type definition for '%s': %w", model.Name, ErrMultipleSymbolDefinition)
			c.Errors = append(c.Errors, err)
			continue
		}

		fields := make(map[string]Field)
		for _, field := range model.Fields {
			_, found := fields[field.Name]
			if found {
				err := fmt.Errorf("duplicate field '%s' in model '%s': %w", field.Name, model.Name, ErrMultipleSymbolDefinition)
				c.Errors = append(c.Errors, err)
				continue
			}
			fields[field.Name] = Field{
				Type: Type{
					Name:    field.TypeName,
					Variant: TypeVariantReference,
				},
				Optional: field.Optional,
			}
		}
		(*typeMap)[model.Name] = Type{
			Name:    model.Name,
			Variant: TypeVariantObject,
			Fields:  fields,
		}
	}
}

func (c *Compiler) resolveTypeReferences(typeMap *map[string]Type) {
	for name, ty := range *typeMap {
		if ty.Variant == TypeVariantReference {
			t, found := (*typeMap)[ty.Name]
			if found {

			}
		}
	}
}

func resolveTypeReference(ty Type, typeMap *map[string]Type) {

}

func (c *Compiler) addBuiltInTypes(m *map[string]Type) {
	(*m)["bool"] = Scalar{
		name: "bool",
	}
	(*m)["int"] = Scalar{
		name: "int",
	}
	(*m)["float"] = Scalar{
		name: "float",
	}
	(*m)["string"] = Scalar{
		name: "string",
	}
}

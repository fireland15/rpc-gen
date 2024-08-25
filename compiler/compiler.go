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
	Fields  map[string]*Field
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
	service.types = c.getDefinedTypes(def)
	service.rpc = c.getRpcs(def)
	return service, nil
}

func (c *Compiler) getRpcs(def parser.ServiceDefinition) map[string]Rpc {
	rpcs := make(map[string]Rpc)

	for _, rpc := range def.Rpc {
		_, exists := rpcs[rpc.Name]
		if exists {
			err := fmt.Errorf("")
			c.Errors = append(c.Errors, err)
		}
	}
}

var ErrMultipleSymbolDefinition = errors.New("symbol defined multiple times")
var ErrUndefinedSymbol = errors.New("the symbol has not been defined")

func (c *Compiler) getDefinedTypes(sd parser.ServiceDefinition) map[string]Type {
	typeMap := make(map[string]Type)

	c.addBuiltInTypes(&typeMap)
	c.addCustomTypes(&sd, &typeMap)
	c.resolveTypeReferences(&typeMap)

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

		fields := make(map[string]*Field)
		for _, field := range model.Fields {
			_, found := fields[field.Name]
			if found {
				err := fmt.Errorf("duplicate field '%s' in model '%s': %w", field.Name, model.Name, ErrMultipleSymbolDefinition)
				c.Errors = append(c.Errors, err)
				continue
			}

			f := new(Field)
			f.Type = Type{
				Name:    field.TypeName,
				Variant: TypeVariantReference,
			}
			f.Optional = field.Optional
			fields[field.Name] = f
		}
		(*typeMap)[model.Name] = Type{
			Name:    model.Name,
			Variant: TypeVariantObject,
			Fields:  fields,
		}
	}
}

func (c *Compiler) resolveTypeReferences(typeMap *map[string]Type) {
	for _, ty := range *typeMap {
		if ty.Variant == TypeVariantObject {
			c.resolveObjectFieldTypeReferences(ty, typeMap)
		}
	}
}

func (c *Compiler) resolveObjectFieldTypeReferences(ty Type, typeMap *map[string]Type) {
	for fieldName, fieldType := range ty.Fields {
		if fieldType.Type.Variant == TypeVariantReference {
			resolvedTy, found := (*typeMap)[fieldType.Type.Name]
			if found {
				fieldType.Type = resolvedTy
			} else {
				err := fmt.Errorf(
					"unknown type '%s' in '%s.%s': %w",
					fieldType.Type.Name,
					ty.Name,
					fieldName,
					ErrUndefinedSymbol)
				c.Errors = append(c.Errors, err)
			}
		}
	}
}

func (c *Compiler) addBuiltInTypes(m *map[string]Type) {
	(*m)["bool"] = Type{
		Name:    "bool",
		Variant: TypeVariantScalar,
	}
	(*m)["int"] = Type{
		Name:    "int",
		Variant: TypeVariantScalar,
	}
	(*m)["float"] = Type{
		Name:    "float",
		Variant: TypeVariantScalar,
	}
	(*m)["string"] = Type{
		Name:    "string",
		Variant: TypeVariantScalar,
	}
}

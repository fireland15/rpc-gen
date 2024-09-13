package compiler

import (
	"errors"
	"fmt"
	"io"

	"github.com/fireland15/rpc-gen/internal/parser"
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
	Rpc   map[string]Rpc
	Types map[string]Type
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
	service.Types = c.getDefinedTypes(def)
	service.Rpc = c.getRpcs(def, &service.Types)

	if len(c.Errors) > 0 {
		err := errors.New("compilation errors")
		for _, e := range c.Errors {
			err = fmt.Errorf("%s:\n\t%w", err, e)
		}
		return nil, err
	}

	return service, nil
}

func (c *Compiler) getRpcs(def parser.ServiceDefinition, typeMap *map[string]Type) map[string]Rpc {
	rpcs := make(map[string]Rpc)

	for _, rpc := range def.Rpc {
		_, exists := rpcs[rpc.Name]
		if exists {
			err := fmt.Errorf("rpc '%s' defined multiple times: %w", rpc.Name, ErrMultipleSymbolDefinition)
			c.Errors = append(c.Errors, err)
		}

		r := Rpc{
			Name: rpc.Name,
		}

		if len(rpc.RequestTypeName) > 0 {
			ty, found := (*typeMap)[rpc.RequestTypeName]
			if !found {
				err := fmt.Errorf("type '%s' is undefined: %w", rpc.RequestTypeName, ErrUndefinedSymbol)
				c.Errors = append(c.Errors, err)
			} else {
				if ty.Variant == TypeVariantScalar {
					err := errors.New("scalar types cannot be directly used in rpc method signatures")
					err = fmt.Errorf("error with rpc method %s: %w", rpc.Name, err)
					c.Errors = append(c.Errors, err)
				} else {
					r.RequestType = &ty
				}
			}
		}

		if len(rpc.ResponseTypeName) > 0 {
			ty, found := (*typeMap)[rpc.ResponseTypeName]
			if !found {
				err := fmt.Errorf("type '%s' is undefined: %w", rpc.ResponseTypeName, ErrUndefinedSymbol)
				c.Errors = append(c.Errors, err)
			} else {
				r.ResponseType = &ty
			}
		}

		rpcs[rpc.Name] = r
	}

	return rpcs
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
	(*m)["date"] = Type{
		Name:    "date",
		Variant: TypeVariantScalar,
	}
	(*m)["uuid"] = Type{
		Name:    "uuid",
		Variant: TypeVariantScalar,
	}
}

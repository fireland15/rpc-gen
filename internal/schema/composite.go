package schema

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

// Composite represents a declared composite (object/model) type in the schema.
//
// A Composite consists of a named set of fields, each with an associated
// type reference and JSON property name.
type Composite interface {
	TypeDef

	// Name returns the logical name of the composite type.
	Name() string

	// Fields returns the fields of the composite, keyed by field name.
	Fields() map[string]Field

	// AddField adds a new field to the composite type.
	AddField(f Field) error
}

// Field represents a single field within a composite type.
type Field interface {
	// Name returns the logical name of the field.
	Name() string

	// JSONPropertyName returns the name of the field as it appears
	// in serialized JSON.
	JSONPropertyName() string

	// Type returns the type reference of the field.
	Type() TypeRef
}

type composite struct {
	name   string
	fields map[string]*field
}

type field struct {
	name             string
	jsonPropertyName string
	ty               TypeRef
}

var _ Composite = (*composite)(nil)
var _ Field = (*field)(nil)
var _ TypeDef = (*composite)(nil)

// NewComposite creates a new composite type with the given name and fields.
//
// Field names and JSON property names must be unique within the composite.
func NewComposite(name string, fields ...Field) (Composite, error) {
	if name == "" {
		return nil, fmt.Errorf("composite name cannot be empty")
	}

	fieldMap := make(map[string]*field, len(fields))
	jsonNames := make(map[string]struct{}, len(fields))

	for _, f := range fields {
		if f.Name() == "" {
			return nil, fmt.Errorf("field name cannot be empty")
		}
		if f.JSONPropertyName() == "" {
			return nil, fmt.Errorf("json property name cannot be empty")
		}
		if f.Type() == nil {
			return nil, fmt.Errorf("field %q has nil type", f.Name())
		}
		if _, ok := fieldMap[f.Name()]; ok {
			return nil, fmt.Errorf("duplicate field name: %q", f.Name())
		}
		if _, ok := jsonNames[f.JSONPropertyName()]; ok {
			return nil, fmt.Errorf(
				"duplicate json property name: %q",
				f.JSONPropertyName(),
			)
		}

		jsonNames[f.JSONPropertyName()] = struct{}{}
		fieldMap[f.Name()] = &field{
			name:             f.Name(),
			jsonPropertyName: f.JSONPropertyName(),
			ty:               f.Type(),
		}
	}

	return &composite{
		name:   name,
		fields: fieldMap,
	}, nil
}

// NewField creates a new field with the given name, JSON property name,
// and type reference.
func NewField(name string, ty TypeRef) Field {
	return &field{
		name:             name,
		jsonPropertyName: strcase.ToLowerCamel(name),
		ty:               ty,
	}
}

// Name implements Composite
func (c *composite) Name() string {
	return c.name
}

// Fields implements Composite
func (c *composite) Fields() map[string]Field {
	out := make(map[string]Field, len(c.fields))
	for k, v := range c.fields {
		out[k] = v
	}
	return out
}

// Kind implements TypeDef
func (c *composite) Kind() TypeKind {
	return TypeKindComposite
}

// AddField adds a new field to the composite type.
//
// Returns an error if the field name or JSON property name is already present,
// or if the field's type is nil.
func (c *composite) AddField(f Field) error {
	if f == nil {
		return fmt.Errorf("field cannot be nil")
	}

	if f.Name() == "" {
		return fmt.Errorf("field name cannot be empty")
	}

	if f.Type() == nil {
		return fmt.Errorf("field %q has nil type", f.Name())
	}

	if _, ok := c.fields[f.Name()]; ok {
		return fmt.Errorf("duplicate field name: %q", f.Name())
	}

	for _, existing := range c.fields {
		if existing.JSONPropertyName() == f.JSONPropertyName() {
			return fmt.Errorf("duplicate JSON property name: %q", f.JSONPropertyName())
		}
	}

	c.fields[f.Name()] = &field{
		name:             f.Name(),
		jsonPropertyName: f.JSONPropertyName(),
		ty:               f.Type(),
	}

	return nil
}

// Name implements Field
func (f *field) Name() string {
	return f.name
}

// JSONPropertyName implements Field
func (f *field) JSONPropertyName() string {
	return f.jsonPropertyName
}

// Type implements Field
func (f *field) Type() TypeRef {
	return f.ty
}

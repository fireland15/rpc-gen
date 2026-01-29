package schema

import (
	"fmt"
	"iter"
	"sort"
)

// Schema is a container for all types and methods defined in an RPC protocol.
// It allows adding and storing TypeDefs (scalars, enums, objects) and Methods.
// Schema ensures that type and method names are unique.
type Schema interface {
	// TypeDefs gets the type definitions in the schema.
	TypeDefs() iter.Seq[TypeDef]

	// TypeDef returns the TypeDef with the name.
	TypeDef(name string) TypeDef

	// Methods returns an iterator over the schema's methods.
	Methods() iter.Seq[Method]

	// Method returns the method with the name.
	Method(name string) Method

	// Items returns an iterator over all items in the schema.
	//
	// Includes methods, types, etc...
	Items() iter.Seq[interface{}]

	// AddType adds a new type definition to the schema.
	// Returns an error if a type with the same name already exists.
	AddType(t TypeDef) error

	// AddMethod adds a new method to the schema.
	// Returns an error if a method with the same name already exists.
	AddMethod(m Method) error

	// Schema converts the builder to a schema
	Schema() Schema
}

type schema struct {
	types   map[string]TypeDef
	methods map[string]Method
}

var _ Schema = (*schema)(nil)

// NewSchema constructs a new empty Schema instance.
func NewSchema() Schema {
	s := &schema{
		types:   make(map[string]TypeDef),
		methods: make(map[string]Method),
	}
	return s
}

// AddType implements Schema.AddType.
//
// Returns an error if a type with the same name is already defined.
func (s *schema) AddType(t TypeDef) error {
	if _, exists := s.types[t.Name()]; exists {
		return fmt.Errorf("type %s already exists", t.Name())
	}
	s.types[t.Name()] = t
	return nil
}

// AddMethod implements Schema.AddMethod.
//
// Returns an error if a method with the same name is already defined.
func (s *schema) AddMethod(m Method) error {
	if _, exists := s.methods[m.Name()]; exists {
		return fmt.Errorf("method %s already exists", m.Name())
	}
	s.methods[m.Name()] = m
	return nil
}

// TypeDefs implements Schema
func (s *schema) TypeDefs() iter.Seq[TypeDef] {
	return func(yield func(TypeDef) bool) {
		names := make([]string, 0, len(s.types))
		for name := range s.types {
			names = append(names, name)
		}

		sort.Strings(names)

		for _, name := range names {
			if !yield(s.types[name]) {
				return
			}
		}
	}
}

// Methods implements Schema
func (s *schema) Methods() iter.Seq[Method] {
	return func(yield func(Method) bool) {
		names := make([]string, 0, len(s.methods))
		for name := range s.methods {
			names = append(names, name)
		}

		sort.Strings(names)

		for _, name := range names {
			if !yield(s.methods[name]) {
				return
			}
		}
	}
}

func (s *schema) Items() iter.Seq[interface{}] {
	return func(yield func(interface{}) bool) {
		for _, t := range s.types {
			if !yield(t) {
				return
			}
		}
		for _, m := range s.methods {
			if !yield(m) {
				return
			}
		}
	}
}

// TypeDef implements Schema
func (s *schema) TypeDef(name string) TypeDef {
	td, ok := s.types[name]
	if !ok {
		return nil
	}
	return td
}

func (s *schema) Method(name string) Method {
	if m, ok := s.methods[name]; ok {
		return m
	}
	return nil
}

func (s *schema) Schema() Schema {
	return s
}

// TypeDef represents a declared type in the schema, which can be a scalar, enum, or composite object.
// Each type must have a unique name within the schema.
type TypeDef interface {
	// Name returns the name of the type.
	Name() string

	// Kind returns the kind of type (scalar, enum, or composite object).
	Kind() TypeKind
}

// TypeKind is an enumeration of the different kinds of schema types.
type TypeKind string

const (
	// TypeKindScalar represents a primitive scalar type (e.g., int, string, boolean).
	TypeKindScalar TypeKind = "scalar"

	// TypeKindEnum represents an enumerated type with a fixed set of values.
	TypeKindEnum TypeKind = "enum"

	// TypeKindComposite represents a composite object type, typically with fields.
	TypeKindComposite TypeKind = "object"
)

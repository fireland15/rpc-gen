package schema

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

// Enumeration represents a declared enum type in the RPC protocol.
type Enumeration interface {
	TypeDef

	// Name returns the logical name of the enumeration.
	Name() string

	// Variants returns the enum variants in declaration order.
	Variants() []EnumVariant
}

// EnumVariant represents a single enumeration variant.
type EnumVariant interface {
	// Name returns the logical name of the variant.
	Name() string

	// SerializedValue returns the wire representation of the variant.
	SerializedValue() string
}

// NewEnumeration creates a new Enumeration with the given name and variant names.
//
// Variant names must be unique. The serialized value of each variant is
// derived from its name.
func NewEnumeration(name string, variantNames ...string) (Enumeration, error) {
	if name == "" {
		return nil, fmt.Errorf("enumeration name cannot be empty")
	}
	if len(variantNames) == 0 {
		return nil, fmt.Errorf("enumeration must have at least one variant")
	}

	seen := make(map[string]struct{}, len(variantNames))
	var variants []*enumVariant

	for _, v := range variantNames {
		if v == "" {
			return nil, fmt.Errorf("enum variant name cannot be empty")
		}
		if _, ok := seen[v]; ok {
			return nil, fmt.Errorf("duplicate enum variant name: %q", v)
		}
		seen[v] = struct{}{}

		variants = append(variants, &enumVariant{
			name:            v,
			serializedValue: deriveEnumSerializedValue(v),
		})
	}

	return &enumeration{
		name:     name,
		variants: variants,
	}, nil
}

type enumeration struct {
	name     string
	variants []*enumVariant
}

type enumVariant struct {
	name            string
	serializedValue string
}

var _ Enumeration = (*enumeration)(nil)
var _ EnumVariant = (*enumVariant)(nil)
var _ TypeDef = (*enumeration)(nil)

func (e *enumeration) Name() string {
	return e.name
}

func (e *enumeration) Variants() []EnumVariant {
	out := make([]EnumVariant, len(e.variants))
	for i, v := range e.variants {
		out[i] = v
	}
	return out
}

func (e *enumeration) Kind() TypeKind {
	return TypeKindEnum
}

func (v *enumVariant) Name() string {
	return v.name
}

func (v *enumVariant) SerializedValue() string {
	return v.serializedValue
}

func deriveEnumSerializedValue(name string) string {
	return strcase.ToScreamingSnake(name)
}

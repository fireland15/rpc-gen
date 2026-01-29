package analysis

import (
	"fmt"

	"github.com/fireland15/rpc-gen/internal/schema"
)

// CheckTypeReferences ensures that all type references in the schema refer to defined types.
//
// Returns nil when all references are valid. Returns an error for the first missing type ref.
func CheckTypeReferences(s schema.Schema) error {
	typeNames := getTypeNames(s)

	for typeDef := range s.TypeDefs() {
		switch t := typeDef.(type) {
		case schema.Composite:
			for _, field := range t.Fields() {
				if err := isValidTypeRef(field.Type(), typeNames); err != nil {
					return fmt.Errorf("invalid typeref: %w", err)
				}
			}
		}
	}

	for m := range s.Methods() {
		if m.ReturnType() != nil {
			if err := isValidTypeRef(m.ReturnType(), typeNames); err != nil {
				return fmt.Errorf("invalid method return type: %w", err)
			}
		}

		for _, arg := range m.Arguments() {
			if err := isValidTypeRef(arg.Type(), typeNames); err != nil {
				return fmt.Errorf("invalid arg type: %w", err)
			}
		}
	}

	return nil
}

func isValidTypeRef(ty schema.TypeRef, definedTypes map[string]struct{}) error {
	switch ty := ty.(type) {
	case schema.NamedTypeRef:
		return isTypeDefined(definedTypes, ty)
	case schema.OptionalTypeRef:
		return isValidTypeRef(ty.InnerType(), definedTypes)
	case schema.ArrayTypeRef:
		return isValidTypeRef(ty.ElementType(), definedTypes)
	default:
		panic("unexpected type ref type.")
	}
}

func getTypeNames(s schema.Schema) map[string]struct{} {
	names := make(map[string]struct{})
	for typeDef := range s.TypeDefs() {
		names[typeDef.Name()] = struct{}{}
	}
	return names
}

func isTypeDefined(definedTypes map[string]struct{}, ty schema.NamedTypeRef) error {
	if _, ok := definedTypes[ty.Name()]; !ok {
		return fmt.Errorf("type %s not defined", ty.Name())
	}
	return nil
}

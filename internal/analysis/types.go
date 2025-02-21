package analysis

import (
	"fmt"
	"slices"

	"github.com/fireland15/rpc-gen/internal/protocol"
)

// Makes sure that type references have a corresponding definition
func CheckTypeReferences(errors *[]string, service *protocol.Protocol) {
	typeNames := getDefinedTypeNames(service.Types)

	for _, m := range service.Types {
		for _, field := range m.Fields {
			if !isTypeDefined(typeNames, field.Type) {
				msg := fmt.Sprintf("undefined type '%s'", field.Type.Name)
				*errors = append(*errors, msg)
			}
		}
	}

	for _, m := range service.Methods {
		for _, p := range m.Parameters {
			if !isTypeDefined(typeNames, p.Type) {
				msg := fmt.Sprintf("undefined type '%s'", p.Type.Name)
				*errors = append(*errors, msg)
			}
		}

		if m.ReturnType != nil && !isTypeDefined(typeNames, m.ReturnType) {
			msg := fmt.Sprintf("undefined type '%s'", m.ReturnType.Name)
			*errors = append(*errors, msg)
		}
	}
}

func isTypeDefined(definedTypes []string, ty *protocol.TypeRef2) bool {
	if !slices.Contains(definedTypes, ty.Name) {
		return false
	} else {
		for idx := range ty.GenericArgs {
			if !isTypeDefined(definedTypes, ty.GenericArgs[idx]) {
				return false
			}
		}
	}

	return true
}

func getDefinedTypeNames(typeDefinitions map[string]*protocol.TypeDefinition) []string {
	names := make([]string, len(typeDefinitions))

	for _, typeDef := range typeDefinitions {
		names = append(names, typeDef.Name)
	}

	names = append(names, "bool", "int", "string", "float", "uuid", "date", "Upload")

	return names
}

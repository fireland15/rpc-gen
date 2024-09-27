package analysis

import (
	"fmt"
	"slices"

	"github.com/fireland15/rpc-gen/internal/model"
)

// Makes sure that type references have a corresponding definition
func CheckTypeReferences(errors *[]string, service model.ServiceDefinition) {
	typeNames := getDefinedTypeNames(service.Models)

	for _, m := range service.Models {
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

		if m.ReturnType != nil && !isTypeDefined(typeNames, *m.ReturnType) {
			msg := fmt.Sprintf("undefined type '%s'", m.ReturnType.Name)
			*errors = append(*errors, msg)
		}
	}
}

func isTypeDefined(definedTypes []string, ty model.Type) bool {
	for ty.Variant != model.TypeVariantNamed {
		if ty.Inner != nil {
			ty = *ty.Inner
		} else {
			panic("non-named types should have an inner type.")
		}
	}

	return slices.Contains(definedTypes, ty.Name)
}

func getDefinedTypeNames(models []model.Model) []string {
	names := make([]string, len(models))

	for _, model := range models {
		names = append(names, model.Name)
	}

	names = append(names, "bool", "int", "string", "float", "uuid", "date")

	return names
}

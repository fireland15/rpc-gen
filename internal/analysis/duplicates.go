package analysis

import (
	"fmt"
	"slices"

	"github.com/fireland15/rpc-gen/internal/model"
)

func CheckForDuplicateModelFields(errors *[]string, service model.ServiceDefinition) {
	for _, m := range service.Models {
		fieldNames := make([]string, 0, len(m.Fields))
		for _, f := range m.Fields {
			if slices.Contains(fieldNames, f.Name) {
				msg := fmt.Sprintf("duplicate field '%s' in model '%s'.", f.Name, m.Name)
				*errors = append(*errors, msg)
			} else {
				fieldNames = append(fieldNames, f.Name)
			}
		}
	}
}

func CheckForDuplicateMethodParameters(errors *[]string, service *model.ServiceDefinition) {
	for _, m := range service.Methods {
		names := make([]string, len(m.Parameters))
		for _, param := range m.Parameters {
			if slices.Contains(names, param.Name) {
				msg := fmt.Sprintf("duplicate parameter \"%s\" in RPC \"%s\"", param.Name, m.Name)
				*errors = append(*errors, msg)
			}
		}
	}
}

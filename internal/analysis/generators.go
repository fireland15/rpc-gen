package analysis

import (
	"fmt"

	"github.com/fireland15/rpc-gen/internal/model"
)

func GenerateMethodParameterModels(service *model.ServiceDefinition) {
	for idx, method := range service.Methods {
		if len(method.Parameters) == 0 {
			continue
		}

		paramsModel := model.Model{
			Name: fmt.Sprintf("%sParams", method.Name),
		}
		for _, param := range method.Parameters {
			paramsModel.Fields = append(paramsModel.Fields, model.Field(param))
		}
		service.Methods[idx].ParameterType.Variant = model.TypeVariantNamed
		service.Methods[idx].ParameterType.Name = paramsModel.Name

		service.Models = append(service.Models, paramsModel)
	}
}

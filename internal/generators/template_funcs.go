package generators

import (
	"text/template"

	"github.com/fireland15/rpc-gen/internal/protocol"
	"github.com/iancoleman/strcase"
)

func addDefaultTemplateFuncs(funcs *template.FuncMap) {
	(*funcs)["toCamel"] = strcase.ToCamel
	(*funcs)["toLowerCamel"] = strcase.ToLowerCamel
	(*funcs)["toSnake"] = strcase.ToSnake

	(*funcs)["hasParameters"] = func(m protocol.MethodDefinition) bool {
		return len(m.Parameters) > 0
	}
}

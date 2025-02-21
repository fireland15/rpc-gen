package typescript

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/fireland15/rpc-gen/internal/protocol"
	"github.com/iancoleman/strcase"
)

type typescriptGenerator struct {
	config   TypescriptConfig
	template *template.Template
}

type TypescriptConfig struct {
	Output       string
	TemplatePath *string
	Types        *map[string]typeConfig `json:"types"`
}

type typeConfig struct {
	// The symbol name to import
	Name string `json:"name"`
	// The module string for the import
	Module string `json:"module"`
	// Used when using * as the import
	Default bool `json:"default"`
	// Used when doing something like `import * as Bananas from "apples"`
	Aliasing string `json:"aliasing"`
}

//go:embed client.template
var clientTemplate string

func NewTypescriptGenerator(config json.RawMessage) (*typescriptGenerator, error) {
	generator := new(typescriptGenerator)

	err := json.Unmarshal(config, &generator.config)
	if err != nil {
		return nil, err
	}

	//funcs := make(template.FuncMap, 0)
	// funcs["resolveType"] = c.resolveType
	// funcs["joinParameters"] = func(m model.Method) string {
	// 	params := make([]string, len(m.Parameters))
	// 	for idx, p := range m.Parameters {
	// 		params[idx] = fmt.Sprintf("%s: %s", strcase.ToLowerCamel(p.Name), c.resolveType(p.Type))
	// 	}
	// 	return strings.Join(params, ", ")
	// }
	// funcs["returnType"] = func(m model.Method) string {
	// 	if m.ReturnType == nil {
	// 		return "void"
	// 	} else {
	// 		return c.resolveType(*m.ReturnType)
	// 	}
	// }
	// tmpl := template.New("ts-client")

	// if generator.config.TemplatePath != nil {
	// 	_, err := tmpl.ParseFiles(*generator.config.TemplatePath)
	// 	if err != nil {
	// 		return generator, err
	// 	}
	// } else {
	// 	_, err := tmpl.Parse(clientTemplate)
	// 	if err != nil {
	// 		return generator, err
	// 	}
	// }

	// generator.template = tmpl

	return generator, nil
}

func (g *typescriptGenerator) Generate(p *protocol.Protocol) error {
	funcs := make(template.FuncMap, 0)

	funcs["formatTypeRef"] = func(typeref *protocol.TypeRef2) string {
		return g.resolveType(p, typeref)
	}
	funcs["defineType"] = func(typeref *protocol.TypeDefinition) bool {
		return !(typeref.Name == "Array" || typeref.Name == "Optional" || typeref.Name == "string")
	}
	funcs["last"] = func(idx int, length int) bool {
		return idx+1 == length
	}
	funcs["lowerCamelCase"] = strcase.ToLowerCamel

	err := os.MkdirAll(filepath.Dir(g.config.Output), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(g.config.Output)
	if err != nil {
		err = fmt.Errorf("problem opening '%s' (GoEchoServerGenerator): %w", g.config.Output, err)
		return err
	}
	defer f.Close()

	templ, err := template.New("ts-client").Funcs(funcs).Parse(clientTemplate)
	if err != nil {
		return fmt.Errorf("parsing error: %w", err)
	}

	err = templ.ExecuteTemplate(f, "ts-client", p)
	if err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	return nil
}

// func mapToTypescript(sd model.ProtocolDefinition) typescriptCodegen {
// 	x := typescriptCodegen{}

// 	// Convert models into typescript types
// 	for _, m := range sd.Models {
// 		tsType := model.Model{
// 			Name:   strcase.ToCamel(m.Name),
// 			Fields: make([]model.Field, len(m.Fields)),
// 		}

// 		for idx, f  := range m.Fields {
// 			tsType.Fields[idx] = model.Field{
// 				Name: strcase.ToLowerCamel(f.Name),

// 			}}
// 		}
// 	}
// 	return x
// }

// func (g *TypescriptGenerator) collectImports(m model.ServiceDefinition) []importStatement {
// 	importStatements := make([]importStatement, 0)
// 	for _, model := range m.Models {
// 		for _, field := range model.Fields {
// 			rootTypeName(field.Type)
// 		}
// 	}
// }

func (g *typescriptGenerator) resolveType(p *protocol.Protocol, typeName *protocol.TypeRef2) string {
	if typeName == nil {
		return "void"
	}
	st, ok := p.GetSerializedType(typeName)
	if !ok {
		return "unknown"
	}
	switch st {
	case protocol.SerializedTypeObject:
		td, exists := p.Types[typeName.Name]
		if !exists {
			return "unknown"
		}
		if td.Name == "Optional" {
			inner := g.resolveType(p, typeName.GenericArgs[0])
			return fmt.Sprintf("%s | null", inner)
		}
		return typeName.Name
	case protocol.SerializedTypeArray:
		inner := g.resolveType(p, typeName.GenericArgs[0])
		return fmt.Sprintf("%s[]", inner)
	default:
		return string(st)
	}
}

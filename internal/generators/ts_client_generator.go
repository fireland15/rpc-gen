package generators

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fireland15/rpc-gen/internal/model"
	"github.com/iancoleman/strcase"
)

type TypescriptClientConfig struct {
	Output string            `json:"output"`
	Types  map[string]string `json:"types"`
}

type TypescriptClientGenerator struct {
	config   TypescriptClientConfig
	template *template.Template
}

//go:embed ts_client.tmpl
var ts_client_template string

func NewTypescriptClientGenerator(config json.RawMessage) (CodeGenerator, error) {
	if config == nil {
		panic("config is nil")
	}

	c := new(TypescriptClientGenerator)

	err := json.Unmarshal(config, &c.config)
	if err != nil {
		return nil, err
	}

	funcs := make(template.FuncMap, 0)
	funcs["toCamel"] = strcase.ToCamel
	funcs["toLowerCamel"] = strcase.ToLowerCamel
	funcs["resolveType"] = c.resolveType
	funcs["joinParameters"] = func(m model.Method) string {
		params := make([]string, len(m.Parameters))
		for idx, p := range m.Parameters {
			params[idx] = fmt.Sprintf("%s: %s", strcase.ToLowerCamel(p.Name), c.resolveType(p.Type))
		}
		return strings.Join(params, ", ")
	}
	funcs["returnType"] = func(m model.Method) string {
		if m.ReturnType == nil {
			return "void"
		} else {
			return c.resolveType(*m.ReturnType)
		}
	}
	funcs["hasParameters"] = func(m model.Method) bool {
		return len(m.Parameters) > 0
	}

	tmpl, err := template.New("ts-client").Funcs(funcs).Parse(ts_client_template)
	if err != nil {
		return nil, err
	}

	c.template = tmpl

	return c, nil
}

func (g *TypescriptClientGenerator) Generate(service *model.ServiceDefinition) error {
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

	_, err = fmt.Fprintln(f, "// This file is autogenerated. Any changes will be overwritten when regenerated.")
	if err != nil {
		return err
	}

	for _, m := range service.Models {
		err = g.template.ExecuteTemplate(f, "model", m)
		if err != nil {
			return err
		}
	}

	_, err = f.WriteString("type Fetcher<P = unknown, R = unknown> = (params: P) => Promise<R>;\n")
	if err != nil {
		return err
	}

	for _, m := range service.Methods {
		err = g.template.ExecuteTemplate(f, "method", m)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *TypescriptClientGenerator) resolveType(typeName model.Type) string {
	if typeName.Variant == model.TypeVariantNamed {
		alias, found := g.config.Types[typeName.Name]
		if !found {
			return typeName.Name
		}
		return alias
	} else if typeName.Variant == model.TypeVariantOptional {
		inner := g.resolveType(*typeName.Inner)
		return fmt.Sprintf("%s | null", inner)
	} else if typeName.Variant == model.TypeVariantArray {
		inner := g.resolveType(*typeName.Inner)
		return fmt.Sprintf("%s[]", inner)
	} else {
		panic("unreachable")
	}
}

package generator

import (
	"embed"
	"fmt"
	"text/template"

	"github.com/fireland15/rpc-gen/internal/schema"
	"github.com/iancoleman/strcase"
)

//go:embed templates/**/*.template
var defaultTemplates embed.FS

type splitRefType struct {
	Type string
	Data any
}

type LanguageConfig struct {
	Scalars map[string]string `json:"scalars"`
	Passes  []PassConfig      `json:"passes"`

	// Inject is a dictionary of strings that can be pulled into templates
	Inject map[string]string `json:"inject"`
}

type Generator struct {
	funcs    template.FuncMap
	language string
	inject   map[string]string
}

func Generate(language string, cfg LanguageConfig, s schema.Schema) error {
	g := &Generator{
		language: language,
	}

	g.funcs = template.FuncMap{
		"toCase": func(s string, c string) string {
			switch c {
			case "camel":
				return strcase.ToLowerCamel(s)
			case "upperCamel":
				return strcase.ToCamel(s)
			case "snake":
				return strcase.ToSnake(s)
			default:
				return s
			}
		},
		"splitTypeRef": func(typeRef schema.TypeRef) splitRefType {
			switch ty := typeRef.(type) {
			case schema.NamedTypeRef:
				return splitRefType{
					Type: "named",
					Data: ty,
				}
			case schema.OptionalTypeRef:
				return splitRefType{
					Type: "optional",
					Data: ty,
				}
			case schema.ArrayTypeRef:
				return splitRefType{
					Type: "array",
					Data: ty,
				}
			}
			panic("unreachable")
		},
		"resolveHostType": func(scalar schema.Scalar) (string, error) {
			n, ok := cfg.Scalars[scalar.Name()]
			if !ok {
				return "", fmt.Errorf("no such scalar: %s", scalar.Name())
			}
			return n, nil
		},
		"resolveHostTypeRef": func(typeRef schema.NamedTypeRef) string {
			n, ok := cfg.Scalars[typeRef.Name()]
			if !ok {
				return typeRef.Name()
			}
			return n
		},
		"lookup": func(key string) (string, error) {
			v, ok := g.inject[key]
			if !ok {
				return "", fmt.Errorf("no such key in Inject: %s", key)
			}
			return v, nil
		},
		"hasDecorator": func(method schema.Method, name string) bool {
			for _, d := range method.Decorators() {
				if d.Name() == name {
					return true
				}
			}
			return false
		},
		"isOptional": func(typeRef schema.TypeRef) bool {
			switch typeRef.(type) {
			case schema.OptionalTypeRef:
				return true
			default:
				return false
			}
		},
	}

	for _, pass := range cfg.Passes {
		inject := cfg.Inject
		for k, v := range pass.Inject {
			inject[k] = v
		}
		g.inject = inject
		if err := g.RenderPass(pass, s); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) getTemplate() (*template.Template, error) {
	t := &template.Template{}
	t = t.Funcs(g.funcs)

	templatePaths := fmt.Sprintf("templates/%s/*.template", g.language)
	t, err := t.ParseFS(defaultTemplates, templatePaths)
	if err != nil {
		return nil, err
	}
	return t, nil
}

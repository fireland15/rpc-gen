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
}

type Generator struct {
	t     *template.Template
	funcs template.FuncMap
}

func Generate(language string, cfg LanguageConfig, s schema.Schema) error {
	t := &template.Template{}

	funcs := template.FuncMap{
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
	}
	t = t.Funcs(funcs)

	templatePaths := fmt.Sprintf("templates/%s/*.template", language)
	t, err := t.ParseFS(defaultTemplates, templatePaths)
	if err != nil {
		return err
	}

	g := &Generator{
		t:     t,
		funcs: funcs,
	}

	for _, pass := range cfg.Passes {
		err = g.RenderPass(pass, s)
		if err != nil {
			return err
		}
	}

	return nil
	//
	//
	//if err := t.ExecuteTemplate(os.Stdout, "preamble", s); err != nil {
	//	return err
	//}
	//
	//for ty := range s.TypeDefs() {
	//	switch ty := ty.(type) {
	//	case schema.Scalar:
	//		if err := t.ExecuteTemplate(os.Stdout, "scalar", ty); err != nil {
	//			return err
	//		}
	//	case schema.Enumeration:
	//		if err := t.ExecuteTemplate(os.Stdout, "enumeration", ty); err != nil {
	//			return err
	//		}
	//	case schema.Composite:
	//		if err := t.ExecuteTemplate(os.Stdout, "composite", ty); err != nil {
	//			return err
	//		}
	//	}
	//}
	//
	//for method := range s.Methods() {
	//	if err := t.ExecuteTemplate(os.Stdout, "method", method); err != nil {
	//		return err
	//	}
	//}

	return nil
}

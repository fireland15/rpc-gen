package generators

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/fireland15/rpc-gen/internal/compiler"
	"github.com/iancoleman/strcase"
)

type GoServerConfig struct {
	Output  string `json:"output"`
	Package string `json:"package"`
	Types   map[string]struct {
		Package   string `json:"package"`
		Namespace string `json:"namespace"`
		TypeName  string `json:"typeName"`
	} `json:"types"`
}

type GoEchoServerGenerator struct {
	config   GoServerConfig
	template *template.Template
}

type goServiceStructField struct {
	Name     string
	Type     string
	JsonName string
}

type goServiceStruct struct {
	Name   string
	Fields []goServiceStructField
}

type goServiceMethod struct {
	Path        string
	Name        string
	RequestType string
	ReturnType  string
}

func (m goServiceMethod) Signature() string {
	params := ""
	if len(m.RequestType) > 0 {
		params = fmt.Sprintf("request *%s", m.RequestType)
	}

	response := "error"
	if len(m.ReturnType) > 0 {
		response = fmt.Sprintf("(%s, error)", m.ReturnType)
	}

	return fmt.Sprintf("%s(%s) %s", m.Name, params, response)
}

func (m goServiceMethod) HasParameters() bool {
	return len(m.RequestType) > 0
}

func (m goServiceMethod) HasResponse() bool {
	return len(m.ReturnType) > 0
}

func (m goServiceMethod) GetPath() string {
	return strcase.ToSnake(m.Name)
}

type goServiceDescriptor struct {
	Package string
	Imports []string
	Structs []goServiceStruct
	Methods []goServiceMethod
}

//go:embed go_echo_server.tmpl
var go_server_template string

func NewGoEchoServerGenerator(config json.RawMessage) (CodeGenerator, error) {
	if config == nil {
		panic("config is nil")
	}

	c := new(GoEchoServerGenerator)

	err := json.Unmarshal(config, &c.config)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("test").Parse(go_server_template)
	if err != nil {
		return nil, err
	}

	c.template = tmpl
	return c, nil
}

func (g *GoEchoServerGenerator) Generate(service *compiler.Service) error {
	desc := goServiceDescriptor{
		Package: g.config.Package,
	}

	for _, t := range g.config.Types {
		desc.Imports = append(desc.Imports, t.Package)
	}

	for _, ty := range service.Types {
		if ty.Variant == compiler.TypeVariantObject {
			s := goServiceStruct{
				Name: ty.Name,
			}
			for name, f := range ty.Fields {
				field := goServiceStructField{
					Name:     strcase.ToCamel(name),
					JsonName: strcase.ToLowerCamel(name),
				}
				typename := f.Type.Name
				importedType, found := g.config.Types[typename]
				if found {
					typename = fmt.Sprintf("%s.%s", importedType.Namespace, importedType.TypeName)
				}
				if f.Optional {
					field.Type = fmt.Sprintf("*%s", typename)
				} else {
					field.Type = typename
				}
				s.Fields = append(s.Fields, field)
			}
			desc.Structs = append(desc.Structs, s)
		}
	}

	for _, rpc := range service.Rpc {
		m := goServiceMethod{
			Name: strcase.ToCamel(rpc.Name),
		}

		if rpc.RequestType != nil {
			m.RequestType = rpc.RequestType.Name
		}

		if rpc.ResponseType != nil {
			m.ReturnType = rpc.ResponseType.Name
		}

		desc.Methods = append(desc.Methods, m)
	}

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

	err = g.template.Execute(f, desc)
	if err != nil {
		err = fmt.Errorf("problem executing template (GoEchoServerGenerator): %w", err)
		return err
	}
	return nil
}

package writer

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"text/template"

	"github.com/fireland15/rpc-gen/internal/compiler"
	"github.com/fireland15/rpc-gen/internal/config"
	"github.com/iancoleman/strcase"
)

type GoServerConfig struct {
	Package string `json:"package"`
	Types   map[string]struct {
		Package   string `json:"package"`
		Namespace string `json:"namespace"`
		TypeName  string `json:"typeName"`
	} `json:"types"`
}

type GoServerWriter struct {
	config   GoServerConfig
	output   io.Writer
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

func NewGoServerWriter(config *config.RpcGenConfig, output io.Writer) (ServerWriter, error) {
	if config == nil {
		panic("config is nil")
	}
	if output == nil {
		panic("output is nil")
	}

	c := new(GoServerWriter)

	rawConfig, found := config.Servers["go"]
	if found {
		err := json.Unmarshal(rawConfig, &c.config)
		if err != nil {
			return nil, err
		}
	}

	c.output = output

	tmpl, err := template.New("test").Parse(go_server_template)
	if err != nil {
		return nil, err
	}

	c.template = tmpl
	return c, nil
}

func (w *GoServerWriter) Write(service *compiler.Service) error {
	desc := goServiceDescriptor{
		Package: w.config.Package,
	}

	for _, t := range w.config.Types {
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
				importedType, found := w.config.Types[typename]
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

	w.template.Execute(w.output, desc)
	return nil
}

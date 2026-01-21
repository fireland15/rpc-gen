package typescript

import (
	"bytes"
	"text/template"

	"github.com/fireland15/rpc-gen/internal/protocol"
)

type TSSDK struct {
	Types   []TSType
	Methods []TSMethod
}

type TSTypeVariant int

const (
	TSTypeVariantScalar TSTypeVariant = iota
	TSTypeVariantObject
	TSTypeVariantArray
	TSTypeVariantUnion
)

type TSType struct {
	Name    string
	Variant TSTypeVariant

	// object
	Fields []TSField

	// array
	ElementType []TSType

	// scalar
	ApiName string

	// union
	Types []TSType
}

type TSField struct {
	Name     string
	Type     string
	Optional bool
}

type TSMethod struct {
	Name       string
	Url        string
	Params     []TSParam
	ReturnType string
	UsesUpload bool
}

type TSParam struct {
	Name string
	Type string
}

func BuildTSSDK(p *protocol.Protocol) *TSSDK {
	sdk := &TSSDK{}

	for _, td := range p.Types {
		if td.Variant != protocol.Compound {
			continue
		}

		t := TSType{Name: td.Name}

		for _, f := range td.Fields {
			t.Fields = append(t.Fields, TSField{
				Name:     f.Name,
				Type:     typeRefToTS(p, f.Type),
				Optional: f.Type.Name == "Optional",
			})
		}

		sdk.Types = append(sdk.Types, t)
	}

	for _, m := range p.Methods {
		method := TSMethod{
			Name:       m.Name,
			Url:        "/" + m.Url(),
			UsesUpload: m.SerializationStrategy() == protocol.SerializationStrategyUpload,
			ReturnType: "void",
		}

		if m.ReturnType != nil {
			method.ReturnType = typeRefToTS(p, m.ReturnType)
		}

		for _, param := range m.Parameters {
			method.Params = append(method.Params, TSParam{
				Name: param.Name,
				Type: typeRefToTS(p, param.Type),
			})
		}

		sdk.Methods = append(sdk.Methods, method)
	}

	return sdk
}

func RenderTSSDK(p *protocol.Protocol) (string, error) {
	sdk := BuildTSSDK(p)

	tmpl, err := template.ParseFiles("./internal/generators/typescript/sdk.ts.template")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, sdk); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func typeRefToTS(p *protocol.Protocol, t *protocol.TypeRef2) string {
	if t == nil {
		return "any"
	}

	switch t.Name {

	case "Array":
		if len(t.GenericArgs) != 1 {
			return "any[]"
		}
		return typeRefToTS(p, t.GenericArgs[0]) + "[]"

	case "Optional":
		if len(t.GenericArgs) != 1 {
			return "any | null"
		}
		return typeRefToTS(p, t.GenericArgs[0]) + " | null"
	}

	// Resolve through protocol
	td, ok := p.Types[t.Name]
	if !ok {
		// Unknown type — assume user-defined object
		return t.Name
	}

	if td.Variant == protocol.Compound {
		return td.Name
	}

	// Scalar — resolve by SerializedType
	if td.SerializedType == nil {
		return "any"
	}

	switch *td.SerializedType {
	case protocol.SerializedTypeString:
		return "string"
	case protocol.SerializedTypeNumber:
		return "number"
	case protocol.SerializedTypeBool:
		return "boolean"
	case protocol.SerializedTypeFile:
		return "Upload"
	default:
		return "any"
	}
}

func collectUsedBrandedTypes(p *protocol.Protocol) map[string]string {
	used := make(map[string]string)

	for _, td := range p.Types {
		if td.Variant != protocol.Scalar {
			continue
		}

		brand, ok := defaultTSBrands[td.Name]
		if !ok {
			continue
		}

		used[td.Name] = brand
	}

	return used
}

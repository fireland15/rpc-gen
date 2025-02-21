package protocol

import "errors"

var ErrDuplicateTypeDefinition = errors.New("duplicate type definition")
var ErrDuplicateMethodDefinition = errors.New("duplicate method definition")

type Protocol struct {
	Name    string
	Methods []*MethodDefinition
	Types   map[string]*TypeDefinition
}

func NewProtocol(name string) *Protocol {
	p := new(Protocol)
	p.Name = name
	p.Methods = make([]*MethodDefinition, 0)
	p.Types = make(map[string]*TypeDefinition)

	st1 := SerializedTypeArray
	p.Types["Array"] = &TypeDefinition{
		Variant:        Compound,
		Name:           "Array",
		Fields:         make([]*FieldDefinition, 0),
		SerializedType: &st1,
	}

	st2 := SerializedTypeObject
	p.Types["Optional"] = &TypeDefinition{
		Variant:        Compound,
		Name:           "Optional",
		Fields:         make([]*FieldDefinition, 0),
		SerializedType: &st2,
	}

	st4 := SerializedTypeString
	p.Types["string"] = &TypeDefinition{
		Variant:        Scalar,
		Name:           "string",
		Fields:         make([]*FieldDefinition, 0),
		SerializedType: &st4,
	}

	st5 := SerializedTypeFile
	p.Types["Upload"] = &TypeDefinition{
		Variant:        Scalar,
		Name:           "Upload",
		Fields:         make([]*FieldDefinition, 0),
		SerializedType: &st5,
	}

	return p
}

func (p *Protocol) AddMethod(method *MethodDefinition) error {
	for idx := range p.Methods {
		if p.Methods[idx].Name == method.Name {
			return ErrDuplicateMethodDefinition
		}
	}
	method.protocol = p
	p.Methods = append(p.Methods, method)
	return nil
}

func (sd *Protocol) AddTypeDefinition(typedef *TypeDefinition) error {
	_, exists := sd.Types[typedef.Name]
	if exists {
		return ErrDuplicateTypeDefinition
	}
	sd.Types[typedef.Name] = typedef
	return nil
}

func (p *Protocol) GetSerializedType(typeref *TypeRef2) (SerializedType, bool) {
	typedef, exists := p.Types[typeref.Name]
	if !exists {
		return "", false
	}

	if typedef.Variant == Compound && typedef.SerializedType == nil {
		return SerializedTypeObject, true
	}

	return *typedef.SerializedType, true
}

func (p *Protocol) MakeType(typeref *TypeRef2) *Type {
	typedef, exists := p.Types[typeref.Name]
	if !exists {
		return nil
	}

	ty := new(Type)
	ty.Name = typedef.Name
	ty.protocol = p
	ty.Fields = make([]*FieldDefinition, 0)

	if typedef.Variant == Scalar {
		return ty
	}

	for idx := range typedef.Fields {
		f := FieldDefinition{
			Name: typedef.Fields[idx].Name,
			Type: typedef.Fields[idx].Type,
		}
		ty.Fields = append(ty.Fields, &f)
	}

	return ty
}

type Type struct {
	Name     string
	Fields   []*FieldDefinition
	protocol *Protocol
}

func (t *Type) ContainsUpload() bool {
	if t.Name == "Upload" {
		return true
	}
	for idx := range t.Fields {
		if t.Fields[idx].Name == "Upload" {
			return true
		} else {
			ty := t.protocol.MakeType(t.Fields[idx].Type)
			if ty == nil {
				continue
			}
			if ty.ContainsUpload() {
				return true
			}
		}
	}

	return false
}

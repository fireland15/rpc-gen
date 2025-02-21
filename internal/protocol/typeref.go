package protocol

import "fmt"

type TypeVariant string

const UploadTypeName = "Upload"

const (
	TypeVariantScalar   TypeVariant = "SCALAR"
	TypeVariantObject   TypeVariant = "OBJECT"
	TypeVariantArray    TypeVariant = "ARRAY"
	TypeVariantOptional TypeVariant = "OPTIONAL"
)

type TypeRef2 struct {
	Name        string
	GenericArgs []*TypeRef2
}

func NewTypeRef(name string) *TypeRef2 {
	return &TypeRef2{
		Name:        name,
		GenericArgs: make([]*TypeRef2, 0),
	}
}

func NewOptionalType(itemType *TypeRef2) *TypeRef2 {
	return &TypeRef2{
		Name: "Optional",
		GenericArgs: []*TypeRef2{
			itemType,
		},
	}
}

func NewArrayType(itemType *TypeRef2) *TypeRef2 {
	return &TypeRef2{
		Name: "Array",
		GenericArgs: []*TypeRef2{
			itemType,
		},
	}
}

type TypeRef struct {
	Name    string
	Variant TypeVariant
	Inner   *TypeRef
}

func (t TypeRef) String() string {
	switch t.Variant {
	case TypeVariantArray:
		return fmt.Sprintf("%s[]", t.Inner.String())
	case TypeVariantOptional:
		return fmt.Sprintf("%s?", t.Inner.String())
	case TypeVariantScalar:
		return t.Name
	case TypeVariantObject:
		return t.Name
	default:
		panic("unreachable")
	}
}

func (t TypeRef) BaseName() string {
	switch t.Variant {
	case TypeVariantArray:
		return t.Inner.BaseName()
	case TypeVariantOptional:
		return t.Inner.BaseName()
	case TypeVariantScalar:
		return t.Name
	case TypeVariantObject:
		return t.Name
	default:
		panic("unreachable")
	}
}

func (t TypeRef) IsUploadType() bool {
	return t.BaseName() == UploadTypeName
}

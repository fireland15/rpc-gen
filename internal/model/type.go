package model

import "fmt"

type TypeVariant int

const (
	TypeVariantNamed TypeVariant = iota
	TypeVariantArray
	TypeVariantOptional
)

type Type struct {
	Name    string
	Variant TypeVariant
	Inner   *Type
}

func (t Type) String() string {
	if t.Variant == TypeVariantNamed {
		return t.Name
	} else if t.Variant == TypeVariantArray {
		return fmt.Sprintf("%s[]", t.Inner.String())
	} else if t.Variant == TypeVariantOptional {
		return fmt.Sprintf("%s?", t.Inner.String())
	}
	panic("unreachable")
}

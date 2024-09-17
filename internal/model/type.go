package model

type TypeVariant int

const (
	TypeVariantReference TypeVariant = iota
	TypeVariantArray
	TypeVariantOptional
)

type Type struct {
	Name    string
	Variant TypeVariant
	Inner   *Type
}

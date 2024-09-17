package model

type Method struct {
	Name       string
	Parameters []MethodParameter
	ReturnType *Type
}

type MethodParameter struct {
	Name string
	Type Type
}

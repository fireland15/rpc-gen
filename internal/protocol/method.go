package protocol

import (
	"errors"

	"github.com/iancoleman/strcase"
)

var ErrDuplicationMethodParameterDefinition = errors.New("duplicate method parameter")

type SerializationStrategy string

const (
	SerializationStrategyDefault SerializationStrategy = "default"
	SerializationStrategyUpload  SerializationStrategy = "upload"
)

type MethodDefinition struct {
	Name                  string
	Parameters            []MethodParameter
	ReturnType            *TypeRef2
	protocol              *Protocol
	serializationStrategy *SerializationStrategy
}

type MethodParameter struct {
	Name string
	Type *TypeRef2
}

func NewMethodDefinition(name string) *MethodDefinition {
	md := new(MethodDefinition)
	md.Name = name
	md.Parameters = make([]MethodParameter, 0)
	return md
}

func (md *MethodDefinition) AddParameter(name string, ty *TypeRef2) error {
	for idx := range md.Parameters {
		if md.Parameters[idx].Name == name {
			return ErrDuplicationMethodParameterDefinition
		}
	}
	md.Parameters = append(
		md.Parameters,
		MethodParameter{
			Name: name,
			Type: ty,
		},
	)
	return nil
}

func (md *MethodDefinition) Url() string {
	return strcase.ToSnake(md.Name)
}

func (md *MethodDefinition) SerializationStrategy() SerializationStrategy {
	for idx := range md.Parameters {
		ty := md.protocol.MakeType(md.Parameters[idx].Type)
		if ty == nil {
			continue
		}
		if ty.ContainsUpload() {
			return SerializationStrategyUpload
		}
	}

	return SerializationStrategyDefault
}

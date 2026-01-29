package schema

import "fmt"

// Scalar represents a declared primitive type in the RPC protocol.
//
// A Scalar defines how a logical type is named within the protocol and how
// its values are serialized when transmitted.
type Scalar interface {
	TypeDef

	// Name returns the logical name of the scalar type.
	Name() string

	// SerializedType returns the JSON primitive type used to serialize values
	// of this scalar.
	SerializedType() JsonType
}

type scalar struct {
	name           string
	serializedType JsonType
}

var _ Scalar = (*scalar)(nil)
var _ TypeDef = (*scalar)(nil)

// NewScalar creates a new Scalar with the given name and serialized type.
//
// It returns an error if the serialized type is not a valid JSON primitive.
func NewScalar(name string, serializedType JsonType) (Scalar, error) {
	if !isValidJsonType(serializedType) {
		return nil, fmt.Errorf("invalid JsonType: %q", serializedType)
	}

	return &scalar{
		name:           name,
		serializedType: serializedType,
	}, nil
}

// Name implements Scalar
func (s *scalar) Name() string {
	return s.name
}

// SerializedType implements Scalar
func (s *scalar) SerializedType() JsonType {
	return s.serializedType
}

// kind implements TypeDef
func (s *scalar) Kind() TypeKind {
	return TypeKindScalar
}

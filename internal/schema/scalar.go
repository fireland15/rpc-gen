package schema

// Scalar represents a declared primitive type in the RPC protocol.
//
// A Scalar defines how a logical type is named within the protocol and how
// its values are serialized when transmitted.
type Scalar interface {
	TypeDef

	// Name returns the logical name of the scalar type.
	Name() string
}

type scalar struct {
	name string
}

var _ Scalar = (*scalar)(nil)
var _ TypeDef = (*scalar)(nil)

// NewScalar creates a new Scalar with the given name and serialized type.
//
// It returns an error if the serialized type is not a valid JSON primitive.
func NewScalar(name string) Scalar {
	return &scalar{
		name: name,
	}
}

// Name implements Scalar
func (s *scalar) Name() string {
	return s.name
}

// kind implements TypeDef
func (s *scalar) Kind() TypeKind {
	return TypeKindScalar
}

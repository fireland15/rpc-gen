package schema

// TypeRef represents a reference to a type in the schema.
//
// TypeRefs are composable and may represent named types as well as
// derived types such as optional or array types.
//
// This interface is sealed and may only be implemented within this package.
type TypeRef interface {
	isTypeRef()
}

type NamedTypeRef interface {
	// Name returns the type's name.
	Name() string
}

// namedTypeRef represents a reference to a named type.
type namedTypeRef struct {
	name string
}

// Name implements NamedTypeRef
func (n *namedTypeRef) Name() string {
	return n.name
}

type OptionalTypeRef interface {
	// InnerType returns a TypeRef for the inner type of the optional
	InnerType() TypeRef
}

// optionalTypeRef represents an optional type (e.g. T?).
type optionalTypeRef struct {
	innerType TypeRef
}

// InnerType implements OptionalTypeRef
func (n *optionalTypeRef) InnerType() TypeRef {
	return n.innerType
}

type ArrayTypeRef interface {
	// ElementType returns a TypeRef for the array's element type.
	ElementType() TypeRef
}

// arrayTypeRef represents an array type (e.g. []T).
type arrayTypeRef struct {
	elementType TypeRef
}

// ElementType implements ArrayTypeRef
func (n *arrayTypeRef) ElementType() TypeRef {
	return n.elementType
}

func (*namedTypeRef) isTypeRef()    {}
func (*optionalTypeRef) isTypeRef() {}
func (*arrayTypeRef) isTypeRef()    {}

var _ TypeRef = (*namedTypeRef)(nil)
var _ TypeRef = (*optionalTypeRef)(nil)
var _ TypeRef = (*arrayTypeRef)(nil)

// NamedType returns a TypeRef that refers to a named type.
//
// The name should correspond to a declared type in the schema.
func NamedType(name string) TypeRef {
	return &namedTypeRef{name: name}
}

// Optional returns a TypeRef representing an optional version of the given type.
//
// For example, Optional(NamedType("User")) represents a nullable User type.
func Optional(inner TypeRef) TypeRef {
	if inner == nil {
		panic("optional type must wrap a non-nil TypeRef")
	}
	return &optionalTypeRef{innerType: inner}
}

// Array returns a TypeRef representing an array of the given type.
//
// For example, Array(NamedType("User")) represents []User.
func Array(element TypeRef) TypeRef {
	if element == nil {
		panic("optional type must wrap a non-nil TypeRef")
	}
	return &arrayTypeRef{elementType: element}
}

// EqualTypeRef returns true when the two TypeRef instances are the same.
func EqualTypeRef(a, b TypeRef) bool {
	if a == b {
		// covers same pointer and both nil
		return true
	}
	if a == nil || b == nil {
		return false
	}

	switch ta := a.(type) {
	case *namedTypeRef:
		tb, ok := b.(*namedTypeRef)
		return ok && ta.name == tb.name

	case *optionalTypeRef:
		tb, ok := b.(*optionalTypeRef)
		return ok && EqualTypeRef(ta.innerType, tb.innerType)

	case *arrayTypeRef:
		tb, ok := b.(*arrayTypeRef)
		return ok && EqualTypeRef(ta.elementType, tb.elementType)

	default:
		// should be unreachable due to sealed interface
		return false
	}
}

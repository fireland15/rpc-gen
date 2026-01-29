package schema_test

import (
	"testing"

	"github.com/fireland15/rpc-gen/internal/schema"
)

func TestNamedType(t *testing.T) {
	name := "User"
	ref := schema.NamedType(name)

	named, ok := ref.(schema.NamedTypeRef)
	if !ok {
		t.Errorf("NamedType() returned wrong type, got %T, want *namedTypeRef", ref)
	}
	if named.Name() != name {
		t.Errorf("NamedType() returned name %q, want %q", named.Name(), name)
	}
}

func TestOptionalType(t *testing.T) {
	inner := schema.NamedType("User")
	ref := schema.Optional(inner)

	opt, ok := ref.(schema.OptionalTypeRef)
	if !ok {
		t.Errorf("Optional() returned wrong type, got %T, want *optionalTypeRef", ref)
	}
	if opt.InnerType() != inner {
		t.Errorf("Optional() returned innerType %v, want %v", opt.InnerType(), inner)
	}
}

func TestOptionalTypeNilPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Optional(nil) did not panic")
		}
	}()
	schema.Optional(nil)
}

func TestArrayType(t *testing.T) {
	element := schema.NamedType("User")
	ref := schema.Array(element)

	arr, ok := ref.(schema.ArrayTypeRef)
	if !ok {
		t.Errorf("Array() returned wrong type, got %T, want *arrayTypeRef", ref)
	}
	if arr.ElementType() != element {
		t.Errorf("Array() returned elementType %v, want %v", arr.ElementType(), element)
	}
}

func TestArrayTypeNilPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Array(nil) did not panic")
		}
	}()
	schema.Array(nil)
}

func TestEqualTypeRef(t *testing.T) {
	tests := []struct {
		name  string
		a     schema.TypeRef
		b     schema.TypeRef
		equal bool
	}{
		{
			name:  "both nil",
			a:     nil,
			b:     nil,
			equal: true,
		},
		{
			name:  "one nil",
			a:     schema.NamedType("User"),
			b:     nil,
			equal: false,
		},
		{
			name:  "same named type",
			a:     schema.NamedType("User"),
			b:     schema.NamedType("User"),
			equal: true,
		},
		{
			name:  "different named types",
			a:     schema.NamedType("User"),
			b:     schema.NamedType("Post"),
			equal: false,
		},
		{
			name:  "optional vs non-optional",
			a:     schema.Optional(schema.NamedType("User")),
			b:     schema.NamedType("User"),
			equal: false,
		},
		{
			name:  "same optional type",
			a:     schema.Optional(schema.NamedType("User")),
			b:     schema.Optional(schema.NamedType("User")),
			equal: true,
		},
		{
			name:  "different optional inner types",
			a:     schema.Optional(schema.NamedType("User")),
			b:     schema.Optional(schema.NamedType("Post")),
			equal: false,
		},
		{
			name:  "same array type",
			a:     schema.Array(schema.NamedType("User")),
			b:     schema.Array(schema.NamedType("User")),
			equal: true,
		},
		{
			name:  "array vs optional",
			a:     schema.Array(schema.NamedType("User")),
			b:     schema.Optional(schema.NamedType("User")),
			equal: false,
		},
		{
			name: "deeply nested equal",
			a: schema.Array(
				schema.Optional(
					schema.Array(
						schema.NamedType("User"),
					),
				),
			),
			b: schema.Array(
				schema.Optional(
					schema.Array(
						schema.NamedType("User"),
					),
				),
			),
			equal: true,
		},
		{
			name: "deeply nested not equal",
			a: schema.Array(
				schema.Optional(
					schema.Array(
						schema.NamedType("User"),
					),
				),
			),
			b: schema.Array(
				schema.Optional(
					schema.Array(
						schema.NamedType("Post"),
					),
				),
			),
			equal: false,
		},
		{
			name:  "same instance shortcut",
			a:     func() schema.TypeRef { t := schema.NamedType("User"); return t }(),
			b:     func() schema.TypeRef { t := schema.NamedType("User"); return t }(),
			equal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.EqualTypeRef(tt.a, tt.b)
			if result != tt.equal {
				t.Fatalf("EqualTypeRef(%#v, %#v) = %v, want %v",
					tt.a, tt.b, result, tt.equal)
			}
		})
	}
}

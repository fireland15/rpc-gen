package parser_test

import (
	"strings"
	"testing"

	"github.com/fireland15/rpc-gen/internal/parser"
	"github.com/fireland15/rpc-gen/internal/schema"
)

func TestParseScalar(t *testing.T) {
	src := `
		scalar UserId = string
		scalar Count = number
	`

	s := mustParse(t, src)

	userId, ok := s.TypeDef("UserId").(schema.Scalar)
	if !ok {
		t.Fatalf("UserId not parsed as scalar")
	}

	if userId.SerializedType() != schema.JsonTypeString {
		t.Errorf("expected UserId to serialize as string")
	}

	count, ok := s.TypeDef("Count").(schema.Scalar)
	if !ok {
		t.Fatalf("Count not parsed as scalar")
	}

	if count.SerializedType() != schema.JsonTypeNumber {
		t.Errorf("expected Count to serialize as number")
	}
}

func TestParseModel(t *testing.T) {
	src := `
		model User {
			id UserId
			name string
			age number?
			tags string[]
		}
	`

	s := mustParse(t, src)

	user, ok := s.TypeDef("User").(schema.Composite)
	if !ok {
		t.Fatalf("User not parsed as model")
	}

	fields := user.Fields()

	if len(fields) != 4 {
		t.Fatalf("expected 4 fields, got %d", len(fields))
	}

	assertField := func(name string, want schema.TypeRef) {
		t.Helper()
		f, ok := fields[name]
		if !ok {
			t.Fatalf("missing field %q", name)
		}
		if !schema.EqualTypeRef(f.Type(), want) {
			t.Errorf("field %q has wrong type: %#v, wanted: %#v", name, f.Type(), want)
		}
	}

	assertField("id", schema.NamedType("UserId"))
	assertField("name", schema.NamedType("string"))
	assertField("age", schema.Optional(schema.NamedType("number")))
	assertField("tags", schema.Array(schema.NamedType("string")))
}

func TestParseEnum(t *testing.T) {
	src := `
		enum Role {
			Admin
			User
			Guest
		}
	`

	s := mustParse(t, src)

	role, ok := s.TypeDef("Role").(schema.Enumeration)
	if !ok {
		t.Fatalf("Role not parsed as enum")
	}

	variants := role.Variants()
	expectedVariantNames := []string{"Admin", "User", "Guest"}
	expectedSerializedValues := []string{"ADMIN", "USER", "GUEST"}

	if len(variants) != len(expectedSerializedValues) {
		t.Fatalf("expected %d variants, got %d", len(expectedSerializedValues), len(variants))
	}

	for i, variant := range expectedVariantNames {
		if variants[i].Name() != variant {
			t.Errorf("expected variant %q, got %q", variant, variants[i].Name())
		}
	}

	for i, v := range expectedSerializedValues {
		if variants[i].SerializedValue() != v {
			t.Errorf("variant %d: expected %q, got %q", i, v, variants[i])
		}
	}
}

func TestParseQuery(t *testing.T) {
	src := `
		query getUser(id UserId) User
	`

	s := mustParse(t, src)

	m := s.Method("getUser")
	if m == nil {
		t.Fatalf("method getUser not found")
	}

	if m.Kind() != schema.MethodKindQuery {
		t.Errorf("expected query method")
	}

	args := m.Arguments()
	if len(args) != 1 {
		t.Fatalf("expected 1 argument, got %d", len(args))
	}

	if args["id"].Type().(schema.NamedTypeRef).Name() != "UserId" {
		t.Errorf("wrong argument type")
	}

	if m.ReturnType().(schema.NamedTypeRef).Name() != "User" {
		t.Errorf("wrong return type")
	}
}

func TestNestedTypes(t *testing.T) {
	src := `
		model Example {
			values string[]?
			matrix number[][] 
		}
	`

	s := mustParse(t, src)

	ex, _ := s.TypeDef("Example").(schema.Composite)

	values, _ := ex.Fields()["values"]
	want := schema.Optional(schema.Array(schema.NamedType("string")))

	if !schema.EqualTypeRef(values.Type(), want) {
		t.Errorf("unexpected type for values")
	}
}

func mustParse(t *testing.T, input string) schema.Schema {
	t.Helper()

	p, err := parser.NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	s, err := p.Parse()
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	return s
}

package parser_test

import (
	"strings"
	"testing"

	"github.com/fireland15/rpc-gen/internal/parser"
	"github.com/fireland15/rpc-gen/internal/schema"
)

func TestParseScalar(t *testing.T) {
	src := `
		scalar UserId
		scalar Count
	`

	s := mustParse(t, src)

	_, ok := s.TypeDef("UserId").(schema.Scalar)
	if !ok {
		t.Fatalf("UserId not parsed as scalar")
	}

	_, ok = s.TypeDef("Count").(schema.Scalar)
	if !ok {
		t.Fatalf("Count not parsed as scalar")
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
		rpc getUser(id UserId) User
	`

	s := mustParse(t, src)

	m := s.Method("getUser")
	if m == nil {
		t.Fatalf("method getUser not found")
	}

	if m.Kind() != schema.MethodKindRpc {
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
func TestParseMethodDecorators(t *testing.T) {
	src := `
		@requires(authenticated, journal_write)
		rpc CreateJournalEntry() JournalEntry
	`

	s := mustParse(t, src)

	m := s.Method("CreateJournalEntry")
	if m == nil {
		t.Fatalf("method CreateJournalEntry not found")
	}

	decorators := m.Decorators()
	if len(decorators) != 1 {
		t.Fatalf("expected 1 decorator, got %d", len(decorators))
	}

	d := decorators[0]

	if d.Name() != "requires" {
		t.Errorf("expected decorator name 'requires', got %q", d.Name())
	}

	expectedArgs := []string{"authenticated", "journal_write"}
	if len(d.Arguments()) != len(expectedArgs) {
		t.Fatalf("expected %d decorator args, got %d", len(expectedArgs), len(d.Arguments()))
	}

	for i, arg := range expectedArgs {
		if d.Arguments()[i] != arg {
			t.Errorf("arg %d: expected %q, got %q", i, arg, d.Arguments()[i])
		}
	}
}

func TestParseMultipleMethodDecorators(t *testing.T) {
	src := `
		@requires(authenticated)
		@rateLimit(user)
		rpc UpdateProfile() User
	`

	s := mustParse(t, src)

	m := s.Method("UpdateProfile")
	if m == nil {
		t.Fatalf("method UpdateProfile not found")
	}

	decorators := m.Decorators()
	if len(decorators) != 2 {
		t.Fatalf("expected 2 decorators, got %d", len(decorators))
	}

	if decorators[0].Name() != "requires" {
		t.Errorf("expected first decorator 'requires', got %q", decorators[0].Name())
	}
	if len(decorators[0].Arguments()) != 1 || decorators[0].Arguments()[0] != "authenticated" {
		t.Errorf("unexpected args for requires decorator: %#v", decorators[0].Arguments())
	}

	if decorators[1].Name() != "rateLimit" {
		t.Errorf("expected second decorator 'rateLimit', got %q", decorators[1].Name())
	}
	if len(decorators[1].Arguments()) != 1 || decorators[1].Arguments()[0] != "user" {
		t.Errorf("unexpected args for rateLimit decorator: %#v", decorators[1].Arguments())
	}
}

func TestParseMethodDecoratorWithoutArgs(t *testing.T) {
	src := `
		@public
		rpc HealthCheck()
	`

	s := mustParse(t, src)

	m := s.Method("HealthCheck")
	if m == nil {
		t.Fatalf("method HealthCheck not found")
	}

	decorators := m.Decorators()
	if len(decorators) != 1 {
		t.Fatalf("expected 1 decorator, got %d", len(decorators))
	}

	d := decorators[0]
	if d.Name() != "public" {
		t.Errorf("expected decorator name 'public', got %q", d.Name())
	}

	if len(d.Arguments()) != 0 {
		t.Errorf("expected no args for public decorator, got %#v", d.Arguments())
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

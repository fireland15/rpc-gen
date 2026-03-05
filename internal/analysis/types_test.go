package analysis_test

import (
	"strings"
	"testing"

	"github.com/fireland15/rpc-gen/internal/analysis"
	"github.com/fireland15/rpc-gen/internal/schema"
)

func TestCheckTypeReferences_AllValid(t *testing.T) {
	s := schema.NewSchema()
	uuid := schema.NewScalar("Uuid")
	s.AddType(uuid)
	colors, _ := schema.NewEnumeration("Colors", "red", "green", "blue")
	s.AddType(colors)
	bird, _ := schema.NewComposite("Bird",
		schema.NewField("id", schema.NamedType("Uuid")),
		schema.NewField("color", schema.Optional(schema.NamedType("Colors"))),
		schema.NewField("homeIds", schema.Array(schema.NamedType("Uuid"))))
	s.AddType(bird)

	arg, _ := schema.NewArgument("id", schema.NamedType("Uuid"))
	m, _ := schema.NewMethod(
		"Exec",
		schema.MethodKindRpc,
		schema.NamedType("Colors"),
		arg)
	s.AddMethod(m)

	if err := analysis.CheckTypeReferences(s); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCheckTypeReferences_MissingFieldType(t *testing.T) {
	s := schema.NewSchema()
	uuid := schema.NewScalar("Uuid")
	s.AddType(uuid)
	colors, _ := schema.NewEnumeration("Colors", "red", "green", "blue")
	s.AddType(colors)
	bird, _ := schema.NewComposite("Bird",
		schema.NewField("id", schema.NamedType("Int")),
		schema.NewField("color", schema.Optional(schema.NamedType("Colors"))),
		schema.NewField("homeIds", schema.Array(schema.NamedType("Uuid"))))
	s.AddType(bird)

	err := analysis.CheckTypeReferences(s)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "Int") {
		t.Fatalf("error should mention missing type, got %v", err)
	}
}

func TestCheckTypeReferences_MissingInnerOfOptionalFieldType(t *testing.T) {
	s := schema.NewSchema()
	uuid := schema.NewScalar("Uuid")
	s.AddType(uuid)
	colors, _ := schema.NewEnumeration("Colors", "red", "green", "blue")
	s.AddType(colors)
	bird, _ := schema.NewComposite("Bird",
		schema.NewField("id", schema.NamedType("Uuid")),
		schema.NewField("color", schema.Optional(schema.NamedType("Bananas"))),
		schema.NewField("homeIds", schema.Array(schema.NamedType("Uuid"))))
	s.AddType(bird)

	err := analysis.CheckTypeReferences(s)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "Bananas") {
		t.Fatalf("error should mention missing type, got %v", err)
	}
}

func TestCheckTypeReferences_MissingMethodReturnType(t *testing.T) {
	s := schema.NewSchema()

	arg, _ := schema.NewArgument("id", schema.NamedType("Uuid"))
	m, _ := schema.NewMethod(
		"Exec",
		schema.MethodKindRpc,
		schema.NamedType("Colors"),
		arg)
	s.AddMethod(m)

	err := analysis.CheckTypeReferences(s)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "Colors") {
		t.Fatalf("error should mention missing type, got %v", err)
	}
}

func TestCheckTypeReferences_MissingMethodArgumentType(t *testing.T) {
	s := schema.NewSchema()

	colors, _ := schema.NewEnumeration("Colors", "red", "green", "blue")
	s.AddType(colors)

	arg, _ := schema.NewArgument("id", schema.NamedType("Uuid"))
	m, _ := schema.NewMethod(
		"Exec",
		schema.MethodKindRpc,
		schema.NamedType("Colors"),
		arg)
	s.AddMethod(m)

	err := analysis.CheckTypeReferences(s)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "Uuid") {
		t.Fatalf("error should mention missing type, got %v", err)
	}
}

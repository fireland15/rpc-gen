package schema_test

import (
	"testing"

	schema2 "github.com/fireland15/rpc-gen/internal/schema"
)

func TestNewSchema(t *testing.T) {
	s := schema2.NewSchema()
	if s == nil {
		t.Fatal("expected schema.NewSchema to return non-nil Schema")
	}
}

func TestAddType_Success(t *testing.T) {
	s := schema2.NewSchema()

	scalar := schema2.NewScalar("String")

	if err := s.AddType(scalar); err != nil {
		t.Fatalf("unexpected error adding type: %v", err)
	}
}

func TestAddType_Duplicate(t *testing.T) {
	s := schema2.NewSchema()

	scalar := schema2.NewScalar("String")
	if err := s.AddType(scalar); err != nil {
		t.Fatalf("unexpected error adding type: %v", err)
	}

	// Add the same type again
	err := s.AddType(scalar)
	if err == nil {
		t.Fatal("expected error when adding duplicate type, got nil")
	}
}

func TestAddMethod_Success(t *testing.T) {
	s := schema2.NewSchema()

	// Create a simple method with no arguments and no return type
	m, err := schema2.NewMethod("GetUser", schema2.MethodKindRpc, nil)
	if err != nil {
		t.Fatalf("unexpected error creating method: %v", err)
	}

	if err := s.AddMethod(m); err != nil {
		t.Fatalf("unexpected error adding method: %v", err)
	}
}

func TestAddMethod_Duplicate(t *testing.T) {
	s := schema2.NewSchema()

	m, _ := schema2.NewMethod("GetUser", schema2.MethodKindRpc, nil)
	if err := s.AddMethod(m); err != nil {
		t.Fatalf("unexpected error adding method: %v", err)
	}

	// Add the same method again
	err := s.AddMethod(m)
	if err == nil {
		t.Fatal("expected error when adding duplicate method, got nil")
	}
}

func TestAddMultipleTypesAndMethods(t *testing.T) {
	s := schema2.NewSchema()

	// Add scalar
	scalar := schema2.NewScalar("Int")
	if err := s.AddType(scalar); err != nil {
		t.Fatalf("unexpected error adding scalar: %v", err)
	}

	// Add enum
	enum, _ := schema2.NewEnumeration("Color", "Red", "Green", "Blue")
	if err := s.AddType(enum); err != nil {
		t.Fatalf("unexpected error adding enum: %v", err)
	}

	// Add methods
	m1, _ := schema2.NewMethod("QueryUser", schema2.MethodKindRpc, nil)
	m2, _ := schema2.NewMethod("UpdateUser", schema2.MethodKindRpc, nil)

	if err := s.AddMethod(m1); err != nil {
		t.Fatalf("unexpected error adding method 1: %v", err)
	}
	if err := s.AddMethod(m2); err != nil {
		t.Fatalf("unexpected error adding method 2: %v", err)
	}
}

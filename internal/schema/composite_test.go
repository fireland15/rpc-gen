package schema_test

import (
	"testing"

	schema2 "github.com/fireland15/rpc-gen/internal/schema"
)

func TestNewField(t *testing.T) {
	ty := schema2.NamedType("User")
	f := schema2.NewField("id", ty)

	if f.Name() != "id" {
		t.Errorf("expected field name 'id', got %q", f.Name())
	}
	if f.JSONPropertyName() != "id" {
		t.Errorf("expected JSON property name 'id', got %q", f.JSONPropertyName())
	}
	if f.Type() != ty {
		t.Errorf("expected type %v, got %v", ty, f.Type())
	}
}

func TestNewComposite(t *testing.T) {
	ty := schema2.NamedType("User")
	f1 := schema2.NewField("id", ty)
	f2 := schema2.NewField("name", ty)

	c, err := schema2.NewComposite("User", f1, f2)
	if err != nil {
		t.Fatalf("unexpected error creating composite: %v", err)
	}

	if c.Name() != "User" {
		t.Errorf("expected composite name 'User', got %q", c.Name())
	}

	fields := c.Fields()
	if len(fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(fields))
	}
	if fields["id"].JSONPropertyName() != "id" {
		t.Errorf("expected JSON property name 'id', got %q", fields["id"].JSONPropertyName())
	}
	if fields["name"].JSONPropertyName() != "name" {
		t.Errorf("expected JSON property name 'name', got %q", fields["name"].JSONPropertyName())
	}
}

func TestNewCompositeErrors(t *testing.T) {
	ty := schema2.NamedType("User")

	// Empty composite name
	_, err := schema2.NewComposite("")
	if err == nil {
		t.Error("expected error for empty composite name")
	}

	// Empty field name
	f1 := schema2.NewField("", ty)
	_, err = schema2.NewComposite("Test", f1)
	if err == nil {
		t.Error("expected error for empty field name")
	}

	// Nil type
	f3 := schema2.NewField("id", nil)
	_, err = schema2.NewComposite("Test", f3)
	if err == nil {
		t.Error("expected error for nil type")
	}

	// Duplicate field names
	f4 := schema2.NewField("id", ty)
	f5 := schema2.NewField("id", ty)
	_, err = schema2.NewComposite("Test", f4, f5)
	if err == nil {
		t.Error("expected error for duplicate field names")
	}
}

func TestAddField(t *testing.T) {
	ty := schema2.NamedType("User")
	c, _ := schema2.NewComposite("User")

	// Add valid field
	f1 := schema2.NewField("id", ty)
	if err := c.AddField(f1); err != nil {
		t.Fatalf("unexpected error adding field: %v", err)
	}

	if len(c.Fields()) != 1 {
		t.Errorf("expected 1 field, got %d", len(c.Fields()))
	}

	// Adding nil field
	if err := c.AddField(nil); err == nil {
		t.Error("expected error when adding nil field")
	}

	// Adding field with empty name
	f2 := schema2.NewField("", ty)
	if err := c.AddField(f2); err == nil {
		t.Error("expected error when adding field with empty name")
	}

	// Adding field with nil type
	f3 := schema2.NewField("name", nil)
	if err := c.AddField(f3); err == nil {
		t.Error("expected error when adding field with nil type")
	}

	// Adding duplicate field name
	f4 := schema2.NewField("id", ty)
	if err := c.AddField(f4); err == nil {
		t.Error("expected error when adding duplicate field name")
	}
}

func TestFieldsImmutability(t *testing.T) {
	ty := schema2.NamedType("User")
	f1 := schema2.NewField("id", ty)
	c, _ := schema2.NewComposite("User", f1)

	fields := c.Fields()
	fields["id"] = schema2.NewField("hack", ty)

	// Original composite should not be affected
	if c.Fields()["id"].Name() != "id" {
		t.Error("fields map should be immutable from external modification")
	}
}

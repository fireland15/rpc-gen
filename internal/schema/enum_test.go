package schema_test

import (
	"reflect"
	"testing"

	"github.com/fireland15/rpc-gen/internal/schema"
)

func TestNewEnumeration_Success(t *testing.T) {
	enum, err := schema.NewEnumeration("Color", "RedRum", "GreenGallop", "BlueBooby")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if enum.Name() != "Color" {
		t.Errorf("Name() = %q, want %q", enum.Name(), "Color")
	}

	variantNames := []string{"RedRum", "GreenGallop", "BlueBooby"}
	variantSerialized := []string{"RED_RUM", "GREEN_GALLOP", "BLUE_BOOBY"}

	variants := enum.Variants()
	if len(variants) != len(variantNames) {
		t.Fatalf("got %d variants, want %d", len(variants), len(variantNames))
	}

	for i, v := range variants {
		if v.Name() != variantNames[i] {
			t.Errorf("variant[%d].Name() = %q, want %q", i, v.Name(), variantNames[i])
		}
		if v.SerializedValue() != variantSerialized[i] {
			t.Errorf(
				"variant[%d].SerializedValue() = %q, want %q",
				i,
				v.SerializedValue(),
				variantSerialized[i],
			)
		}
	}
}

func TestNewEnumeration_DuplicateVariants(t *testing.T) {
	_, err := schema.NewEnumeration("Color", "RED", "GREEN", "RED")
	if err == nil {
		t.Fatal("expected error for duplicate variant names, got nil")
	}
}

func TestNewEnumeration_EmptyEnumName(t *testing.T) {
	_, err := schema.NewEnumeration("", "RED")
	if err == nil {
		t.Fatal("expected error for empty enumeration name, got nil")
	}
}

func TestNewEnumeration_NoVariants(t *testing.T) {
	_, err := schema.NewEnumeration("Color")
	if err == nil {
		t.Fatal("expected error for empty variant list, got nil")
	}
}

func TestNewEnumeration_EmptyVariantName(t *testing.T) {
	_, err := schema.NewEnumeration("Color", "RED", "")
	if err == nil {
		t.Fatal("expected error for empty variant name, got nil")
	}
}

func TestEnumeration_VariantsAreImmutable(t *testing.T) {
	enum, err := schema.NewEnumeration("Color", "RED", "GREEN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	variants := enum.Variants()

	// Mutate returned slice
	variants[0] = nil

	// Fetch again — should be unaffected
	variants2 := enum.Variants()
	if variants2[0] == nil {
		t.Fatal("Variants() exposes internal slice; mutation leaked")
	}
}

func TestEnumeration_VariantOrderPreserved(t *testing.T) {
	input := []string{"FIRST", "SECOND", "THIRD"}

	enum, err := schema.NewEnumeration("OrderTest", input...)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got []string
	for _, v := range enum.Variants() {
		got = append(got, v.Name())
	}

	if !reflect.DeepEqual(got, input) {
		t.Errorf("variant order = %v, want %v", got, input)
	}
}

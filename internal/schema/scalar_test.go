package schema_test

import (
	"testing"

	schema2 "github.com/fireland15/rpc-gen/internal/schema"
)

func TestNewScalar(t *testing.T) {
	scalar, err := schema2.NewScalar("steve", schema2.JsonTypeString)

	if err != nil {
		t.Fatal(err)
	}
	if scalar.Name() != "steve" {
		t.Errorf("name should be steve")
	}
	if scalar.SerializedType() != schema2.JsonTypeString {
		t.Errorf("serialized type should be %s", schema2.JsonTypeString)
	}
}

func TestNewScalarReturnsError(t *testing.T) {
	scalar, err := schema2.NewScalar("steve", "lol-nope")

	if err == nil {
		t.Fatal("Expected an error")
	}
	if scalar != nil {
		{
			t.Fatal("Expected a nil result")
		}
	}
}

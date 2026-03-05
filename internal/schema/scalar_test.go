package schema_test

import (
	"testing"

	schema2 "github.com/fireland15/rpc-gen/internal/schema"
)

func TestNewScalar(t *testing.T) {
	scalar := schema2.NewScalar("steve")

	if scalar.Name() != "steve" {
		t.Errorf("name should be steve")
	}
}

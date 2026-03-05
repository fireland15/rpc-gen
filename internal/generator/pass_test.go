package generator

import (
	"slices"
	"testing"

	"github.com/fireland15/rpc-gen/internal/schema"
)

func TestPass(t *testing.T) {
	s := schema.NewSchema()
	e1, _ := schema.NewEnumeration("Colors", "A", "B", "c")
	e2, _ := schema.NewEnumeration("Colors2", "A", "B", "c")
	_ = s.AddType(e1)
	_ = s.AddType(e2)

	filteredItems := slices.Collect(filtered([]string{"enum"}, s.Schema().Items()))
	if len(filteredItems) != 2 {
		t.Errorf("Expected 2 items, got %d", len(filteredItems))
	}
}

func TestPassMultiFilter(t *testing.T) {
	s := schema.NewSchema()
	e1, _ := schema.NewEnumeration("Colors", "A", "B", "c")
	e2, _ := schema.NewEnumeration("Colors2", "A", "B", "c")
	_ = s.AddType(e1)
	_ = s.AddType(e2)

	s1 := schema.NewScalar("Scalar")
	_ = s.AddType(s1)

	m, _ := schema.NewMethod("method", schema.MethodKindRpc, nil)
	_ = s.AddMethod(m)

	filteredItems := slices.Collect(filtered([]string{"enum", "scalar"}, s.Schema().Items()))
	if len(filteredItems) != 3 {
		t.Errorf("Expected 3 items, got %d", len(filteredItems))
	}
}

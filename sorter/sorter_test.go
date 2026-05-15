package sorter_test

import (
	"testing"

	"github.com/envlock/schema"
	"github.com/envlock/sorter"
)

func makeSchema() schema.Schema {
	return schema.Schema{
		Vars: []schema.EnvVar{
			{Key: "ZEBRA", Required: false, Description: "last alpha"},
			{Key: "ALPHA", Required: true, Description: "first alpha"},
			{Key: "MIDDLE", Required: false, Description: "middle entry"},
			{Key: "BETA", Required: true, Description: "another required"},
		},
	}
}

func TestSort_ByKey(t *testing.T) {
	s := makeSchema()
	result := sorter.Sort(s, sorter.SortByKey)

	expected := []string{"ALPHA", "BETA", "MIDDLE", "ZEBRA"}
	for i, ev := range result.Vars {
		if ev.Key != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], ev.Key)
		}
	}
}

func TestSort_ByRequired(t *testing.T) {
	s := makeSchema()
	result := sorter.Sort(s, sorter.SortByRequired)

	// First two should be required
	for i := 0; i < 2; i++ {
		if !result.Vars[i].Required {
			t.Errorf("index %d should be required, got key %q", i, result.Vars[i].Key)
		}
	}
	// Last two should be optional
	for i := 2; i < 4; i++ {
		if result.Vars[i].Required {
			t.Errorf("index %d should be optional, got key %q", i, result.Vars[i].Key)
		}
	}
}

func TestSort_ByDescription(t *testing.T) {
	s := makeSchema()
	result := sorter.Sort(s, sorter.SortByDescription)

	for i := 1; i < len(result.Vars); i++ {
		if result.Vars[i].Description < result.Vars[i-1].Description {
			t.Errorf("index %d (%q) should come after %d (%q)",
				i, result.Vars[i].Description, i-1, result.Vars[i-1].Description)
		}
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	s := makeSchema()
	originalFirst := s.Vars[0].Key

	sorter.Sort(s, sorter.SortByKey)

	if s.Vars[0].Key != originalFirst {
		t.Errorf("original schema was mutated: expected first key %q, got %q", originalFirst, s.Vars[0].Key)
	}
}

func TestSort_DefaultFallsBackToKey(t *testing.T) {
	s := makeSchema()
	result := sorter.Sort(s, sorter.SortOrder("unknown"))

	expected := []string{"ALPHA", "BETA", "MIDDLE", "ZEBRA"}
	for i, ev := range result.Vars {
		if ev.Key != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], ev.Key)
		}
	}
}

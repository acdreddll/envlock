package stripper_test

import (
	"testing"

	"github.com/your-org/envlock/schema"
	"github.com/your-org/envlock/stripper"
)

func makeSchema() schema.Schema {
	return schema.Schema{
		{Key: "DB_URL", Description: "Database URL", Default: "localhost", Required: true, Sensitive: true, Group: "db", Tags: []string{"core"}, Pattern: "^postgres"},
		{Key: "API_KEY", Description: "API key", Sensitive: true, Required: true, Deprecated: true, RemoveBy: "2025-01-01"},
		{Key: "LOG_LEVEL", Description: "Log verbosity", Default: "info", Group: "logging"},
	}
}

func TestStrip_Description(t *testing.T) {
	s := makeSchema()
	out, err := stripper.Strip(s, stripper.Options{Fields: []string{"description"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, v := range out {
		if v.Description != "" {
			t.Errorf("expected empty description for %s, got %q", v.Key, v.Description)
		}
	}
}

func TestStrip_Default(t *testing.T) {
	s := makeSchema()
	out, err := stripper.Strip(s, stripper.Options{Fields: []string{"default"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, v := range out {
		if v.Default != "" {
			t.Errorf("expected empty default for %s, got %q", v.Key, v.Default)
		}
	}
}

func TestStrip_MultipleFields(t *testing.T) {
	s := makeSchema()
	out, err := stripper.Strip(s, stripper.Options{Fields: []string{"sensitive", "group", "tags"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, v := range out {
		if v.Sensitive {
			t.Errorf("%s: sensitive should be false", v.Key)
		}
		if v.Group != "" {
			t.Errorf("%s: group should be empty", v.Key)
		}
		if len(v.Tags) != 0 {
			t.Errorf("%s: tags should be nil/empty", v.Key)
		}
	}
}

func TestStrip_DeprecatedAndRemoveBy(t *testing.T) {
	s := makeSchema()
	out, err := stripper.Strip(s, stripper.Options{Fields: []string{"deprecated", "remove_by"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, v := range out {
		if v.Deprecated {
			t.Errorf("%s: deprecated should be false", v.Key)
		}
		if v.RemoveBy != "" {
			t.Errorf("%s: remove_by should be empty", v.Key)
		}
	}
}

func TestStrip_DoesNotMutateOriginal(t *testing.T) {
	s := makeSchema()
	origDesc := s[0].Description
	_, err := stripper.Strip(s, stripper.Options{Fields: []string{"description"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s[0].Description != origDesc {
		t.Error("original schema was mutated")
	}
}

func TestStrip_UnknownField(t *testing.T) {
	s := makeSchema()
	_, err := stripper.Strip(s, stripper.Options{Fields: []string{"nonexistent"}})
	if err == nil {
		t.Error("expected error for unknown field, got nil")
	}
}

func TestStrip_EmptyFields(t *testing.T) {
	s := makeSchema()
	out, err := stripper.Strip(s, stripper.Options{Fields: []string{}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(s) {
		t.Errorf("expected %d vars, got %d", len(s), len(out))
	}
}

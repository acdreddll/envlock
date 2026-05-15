package grouper_test

import (
	"testing"

	"github.com/envlock/grouper"
	"github.com/envlock/schema"
)

func makeSchema(envs []schema.EnvVar) schema.Schema {
	return schema.Schema{Envs: envs}
}

func TestGroupBy_Required(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		{Key: "A", Required: true},
		{Key: "B", Required: false},
		{Key: "C", Required: true},
	})

	result := grouper.GroupSchema(s, grouper.GroupByRequired)

	if len(result.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(result.Groups))
	}
	if result.Groups[0].Name != "optional" {
		t.Errorf("expected first group 'optional', got %q", result.Groups[0].Name)
	}
	if result.Groups[1].Name != "required" {
		t.Errorf("expected second group 'required', got %q", result.Groups[1].Name)
	}
	if len(result.Groups[1].Entries) != 2 {
		t.Errorf("expected 2 required entries, got %d", len(result.Groups[1].Entries))
	}
}

func TestGroupBy_Sensitive(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		{Key: "SECRET", Sensitive: true},
		{Key: "HOST", Sensitive: false},
		{Key: "TOKEN", Sensitive: true},
	})

	result := grouper.GroupSchema(s, grouper.GroupBySensitive)

	if len(result.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(result.Groups))
	}
	if result.Groups[0].Name != "non-sensitive" {
		t.Errorf("expected 'non-sensitive', got %q", result.Groups[0].Name)
	}
	if result.Groups[1].Name != "sensitive" {
		t.Errorf("expected 'sensitive', got %q", result.Groups[1].Name)
	}
}

func TestGroupBy_Group(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		{Key: "DB_HOST", Group: "database"},
		{Key: "DB_PORT", Group: "database"},
		{Key: "APP_PORT", Group: "app"},
		{Key: "ORPHAN"},
	})

	result := grouper.GroupSchema(s, grouper.GroupByGroup)

	if len(result.Groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(result.Groups))
	}
	if result.Groups[0].Name != "app" {
		t.Errorf("expected 'app', got %q", result.Groups[0].Name)
	}
	if result.Groups[2].Name != "ungrouped" {
		t.Errorf("expected 'ungrouped', got %q", result.Groups[2].Name)
	}
}

func TestGroupBy_EmptySchema(t *testing.T) {
	s := makeSchema([]schema.EnvVar{})
	result := grouper.GroupSchema(s, grouper.GroupByRequired)
	if len(result.Groups) != 0 {
		t.Errorf("expected no groups for empty schema, got %d", len(result.Groups))
	}
}

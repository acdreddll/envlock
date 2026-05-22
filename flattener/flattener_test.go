package flattener_test

import (
	"testing"

	"github.com/envlock/flattener"
	"github.com/envlock/schema"
)

func makeSchema(vars []schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key string, required, sensitive bool, def, desc string) schema.EnvVar {
	return schema.EnvVar{
		Key:         key,
		Required:    required,
		Sensitive:   sensitive,
		Default:     def,
		Description: desc,
	}
}

func TestFlatten_AllPresent(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		ev("APP_ENV", true, false, "", "Environment name"),
		ev("PORT", false, false, "8080", "HTTP port"),
	})
	env := map[string]string{"APP_ENV": "production", "PORT": "9090"}

	result, err := flattener.Flatten(s, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result.Entries))
	}
	m := result.ToMap()
	if m["PORT"] != "9090" {
		t.Errorf("expected PORT=9090, got %s", m["PORT"])
	}
}

func TestFlatten_FallsBackToDefault(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		ev("PORT", false, false, "8080", "HTTP port"),
	})
	result, err := flattener.Flatten(s, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ToMap()["PORT"] != "8080" {
		t.Errorf("expected default PORT=8080")
	}
}

func TestFlatten_MissingRequired(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		ev("SECRET_KEY", true, true, "", "Secret signing key"),
	})
	_, err := flattener.Flatten(s, map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing required key, got nil")
	}
}

func TestFlatten_OptionalMissingIsEmpty(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		ev("OPTIONAL_VAR", false, false, "", "Optional"),
	})
	result, err := flattener.Flatten(s, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := result.ToMap()["OPTIONAL_VAR"]; v != "" {
		t.Errorf("expected empty string, got %q", v)
	}
}

func TestFlatten_SortedByKey(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		ev("ZEBRA", false, false, "z", ""),
		ev("ALPHA", false, false, "a", ""),
		ev("MIDDLE", false, false, "m", ""),
	})
	result, err := flattener.Flatten(s, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := []string{result.Entries[0].Key, result.Entries[1].Key, result.Entries[2].Key}
	expected := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("position %d: expected %s, got %s", i, expected[i], k)
		}
	}
}

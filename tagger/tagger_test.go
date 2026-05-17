package tagger

import (
	"testing"

	"github.com/user/envlock/schema"
)

func makeSchema(vars []schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key string, tags []string) schema.EnvVar {
	return schema.EnvVar{Key: key, Tags: tags}
}

func TestTag_BuildsIndex(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		ev("DB_HOST", []string{"database", "network"}),
		ev("DB_PORT", []string{"database"}),
		ev("API_KEY", []string{"auth"}),
	})

	r, err := Tag(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(r.Index["database"]) != 2 {
		t.Errorf("expected 2 keys under 'database', got %d", len(r.Index["database"]))
	}
	if len(r.Index["network"]) != 1 || r.Index["network"][0] != "DB_HOST" {
		t.Errorf("expected DB_HOST under 'network'")
	}
	if len(r.Index["auth"]) != 1 || r.Index["auth"][0] != "API_KEY" {
		t.Errorf("expected API_KEY under 'auth'")
	}
}

func TestTag_Untagged(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		ev("PLAIN_VAR", nil),
		ev("ANOTHER", nil),
		ev("TAGGED", []string{"infra"}),
	})

	r, err := Tag(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(r.Untagged) != 2 {
		t.Errorf("expected 2 untagged vars, got %d", len(r.Untagged))
	}
	if r.Untagged[0] != "ANOTHER" || r.Untagged[1] != "PLAIN_VAR" {
		t.Errorf("untagged not sorted correctly: %v", r.Untagged)
	}
}

func TestTag_EmptySchema(t *testing.T) {
	r, err := Tag(makeSchema(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Index) != 0 {
		t.Errorf("expected empty index for empty schema")
	}
}

func TestTag_EmptyTagError(t *testing.T) {
	s := makeSchema([]schema.EnvVar{
		ev("BAD_VAR", []string{""}),
	})
	_, err := Tag(s)
	if err == nil {
		t.Error("expected error for empty tag string, got nil")
	}
}

func TestKeysForTag_Missing(t *testing.T) {
	r := Result{Index: make(TagIndex)}
	keys := KeysForTag(r, "nonexistent")
	if len(keys) != 0 {
		t.Errorf("expected empty slice for missing tag, got %v", keys)
	}
}

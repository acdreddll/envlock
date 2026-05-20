package freezer_test

import (
	"testing"

	"github.com/user/envlock/freezer"
	"github.com/user/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key string, opts ...func(*schema.EnvVar)) schema.EnvVar {
	v := schema.EnvVar{Key: key}
	for _, o := range opts {
		o(&v)
	}
	return v
}

func required(v *schema.EnvVar)  { v.Required = true }
func sensitive(v *schema.EnvVar) { v.Sensitive = true }
func withDefault(d string) func(*schema.EnvVar) {
	return func(v *schema.EnvVar) { v.Default = d }
}

func TestFreeze_BasicResolution(t *testing.T) {
	s := makeSchema(ev("APP_HOST", required), ev("APP_PORT", withDefault("8080")))
	env := map[string]string{"APP_HOST": "localhost"}

	snap, err := freezer.Freeze(s, env, "envlock.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap.Entries))
	}
}

func TestFreeze_MissingRequired(t *testing.T) {
	s := makeSchema(ev("DB_URL", required))
	_, err := freezer.Freeze(s, map[string]string{}, "envlock.yaml")
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
}

func TestFreeze_OptionalMissingOmitted(t *testing.T) {
	s := makeSchema(ev("OPTIONAL_KEY"))
	snap, err := freezer.Freeze(s, map[string]string{}, "envlock.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(snap.Entries))
	}
}

func TestFreeze_DefaultMarked(t *testing.T) {
	s := makeSchema(ev("TIMEOUT", withDefault("30s")))
	snap, err := freezer.Freeze(s, map[string]string{}, "envlock.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !snap.Entries[0].FromDefault {
		t.Error("expected FromDefault=true")
	}
}

func TestFreeze_Redacted(t *testing.T) {
	s := makeSchema(ev("SECRET", sensitive))
	env := map[string]string{"SECRET": "super-secret"}
	snap, _ := freezer.Freeze(s, env, "envlock.yaml")

	redacted := snap.Redacted()
	if redacted.Entries[0].Value != "***" {
		t.Errorf("expected masked value, got %q", redacted.Entries[0].Value)
	}
	// original unchanged
	if snap.Entries[0].Value != "super-secret" {
		t.Error("original snapshot should not be mutated")
	}
}

func TestFreeze_EntriesSorted(t *testing.T) {
	s := makeSchema(ev("ZEBRA", required), ev("ALPHA", required))
	env := map[string]string{"ZEBRA": "z", "ALPHA": "a"}
	snap, _ := freezer.Freeze(s, env, "envlock.yaml")

	if snap.Entries[0].Key != "ALPHA" || snap.Entries[1].Key != "ZEBRA" {
		t.Error("entries should be sorted alphabetically by key")
	}
}

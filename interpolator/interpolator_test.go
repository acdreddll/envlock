package interpolator_test

import (
	"testing"

	"github.com/envlock/interpolator"
	"github.com/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key, def string) schema.EnvVar {
	return schema.EnvVar{Key: key, Default: def}
}

func TestInterpolate_NoReferences(t *testing.T) {
	s := makeSchema(ev("HOST", "localhost"), ev("PORT", "8080"))
	env := map[string]string{"HOST": "example.com", "PORT": "9090"}

	results, err := interpolator.Interpolate(s, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Expanded {
			t.Errorf("expected no expansion for %s", r.Key)
		}
	}
}

func TestInterpolate_ExpandsReference(t *testing.T) {
	s := makeSchema(ev("BASE_URL", "http://${HOST}:${PORT}"))
	env := map[string]string{
		"HOST":     "api.example.com",
		"PORT":     "443",
		"BASE_URL": "http://${HOST}:${PORT}",
	}

	results, err := interpolator.Interpolate(s, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if !r.Expanded {
		t.Error("expected BASE_URL to be expanded")
	}
	want := "http://api.example.com:443"
	if r.Resolved != want {
		t.Errorf("got %q, want %q", r.Resolved, want)
	}
	if env["BASE_URL"] != want {
		t.Errorf("env map not updated: got %q", env["BASE_URL"])
	}
}

func TestInterpolate_UndefinedReference(t *testing.T) {
	s := makeSchema(ev("DSN", "postgres://${DB_HOST}/mydb"))
	env := map[string]string{"DSN": "postgres://${DB_HOST}/mydb"}

	_, err := interpolator.Interpolate(s, env)
	if err == nil {
		t.Fatal("expected error for undefined reference, got nil")
	}
}

func TestInterpolate_SkipsMissingKeys(t *testing.T) {
	s := makeSchema(ev("DEFINED", ""), ev("ABSENT", ""))
	env := map[string]string{"DEFINED": "value"}

	results, err := interpolator.Interpolate(s, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "DEFINED" {
		t.Errorf("unexpected key: %s", results[0].Key)
	}
}

func TestInterpolate_UsesDefaultWhenEnvMissing(t *testing.T) {
	s := makeSchema(ev("TIMEOUT", "30s"))
	// TIMEOUT is not set in env; interpolator should fall back to the schema default.
	env := map[string]string{}

	results, err := interpolator.Interpolate(s, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Key != "TIMEOUT" {
		t.Errorf("unexpected key: %s", r.Key)
	}
	if r.Resolved != "30s" {
		t.Errorf("expected default value %q, got %q", "30s", r.Resolved)
	}
}

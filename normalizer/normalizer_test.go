package normalizer_test

import (
	"testing"

	"github.com/envlock/normalizer"
	"github.com/envlock/schema"
)

func makeSchema(keys ...string) schema.Schema {
	vars := make([]schema.EnvVar, len(keys))
	for i, k := range keys {
		vars[i] = schema.EnvVar{Key: k}
	}
	return schema.Schema{Vars: vars}
}

func TestNormalize_TrimSpace(t *testing.T) {
	s := makeSchema("HOST", "PORT")
	env := map[string]string{"HOST": "  localhost  ", "PORT": "8080"}
	opts := normalizer.DefaultOptions()

	results := normalizer.Normalize(s, env, opts)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Key == "HOST" {
			if r.Normalized != "localhost" {
				t.Errorf("HOST: expected 'localhost', got %q", r.Normalized)
			}
			if !r.Changed {
				t.Error("HOST: expected Changed=true")
			}
		}
		if r.Key == "PORT" && r.Changed {
			t.Error("PORT: expected Changed=false")
		}
	}
}

func TestNormalize_UpperCase(t *testing.T) {
	s := makeSchema("MODE")
	env := map[string]string{"MODE": "production"}
	opts := normalizer.Options{TrimSpace: false, UpperCase: true}

	results := normalizer.Normalize(s, env, opts)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Normalized != "PRODUCTION" {
		t.Errorf("expected 'PRODUCTION', got %q", results[0].Normalized)
	}
}

func TestNormalize_CollapseWhitespace(t *testing.T) {
	s := makeSchema("DESCRIPTION")
	env := map[string]string{"DESCRIPTION": "hello   world  foo"}
	opts := normalizer.Options{TrimSpace: true, CollapseWS: true}

	results := normalizer.Normalize(s, env, opts)

	if results[0].Normalized != "hello world foo" {
		t.Errorf("expected 'hello world foo', got %q", results[0].Normalized)
	}
}

func TestNormalize_SkipsUnknownKeys(t *testing.T) {
	s := makeSchema("KNOWN")
	env := map[string]string{"KNOWN": "val", "UNKNOWN": "  extra  "}
	opts := normalizer.DefaultOptions()

	results := normalizer.Normalize(s, env, opts)

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "KNOWN" {
		t.Errorf("expected key KNOWN, got %s", results[0].Key)
	}
}

func TestNormalize_MissingEnvKey(t *testing.T) {
	s := makeSchema("MISSING")
	env := map[string]string{}
	opts := normalizer.DefaultOptions()

	results := normalizer.Normalize(s, env, opts)

	if len(results) != 0 {
		t.Errorf("expected 0 results for missing key, got %d", len(results))
	}
}

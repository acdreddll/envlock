package masker_test

import (
	"testing"

	"github.com/envlock/masker"
	"github.com/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key string, sensitive bool) schema.EnvVar {
	return schema.EnvVar{Key: key, Sensitive: sensitive}
}

func TestMask_SensitiveOnlyMasksSensitive(t *testing.T) {
	s := makeSchema(ev("API_KEY", true), ev("APP_ENV", false))
	env := map[string]string{"API_KEY": "supersecret", "APP_ENV": "production"}
	opts := masker.DefaultOptions()

	results := masker.Mask(s, env, opts)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Key == "API_KEY" {
			if !r.WasMasked {
				t.Error("expected API_KEY to be masked")
			}
			if r.Masked == "supersecret" {
				t.Error("expected API_KEY value to be obscured")
			}
		}
		if r.Key == "APP_ENV" {
			if r.WasMasked {
				t.Error("expected APP_ENV not to be masked")
			}
			if r.Masked != "production" {
				t.Errorf("expected APP_ENV value unchanged, got %q", r.Masked)
			}
		}
	}
}

func TestMask_AllVarsWhenSensitiveOnlyFalse(t *testing.T) {
	s := makeSchema(ev("DB_URL", false), ev("TOKEN", true))
	env := map[string]string{"DB_URL": "postgres://localhost", "TOKEN": "abc123"}
	opts := masker.DefaultOptions()
	opts.SensitiveOnly = false

	results := masker.Mask(s, env, opts)
	for _, r := range results {
		if !r.WasMasked {
			t.Errorf("expected %s to be masked", r.Key)
		}
	}
}

func TestMask_VisiblePrefixSuffix(t *testing.T) {
	s := makeSchema(ev("SECRET", true))
	env := map[string]string{"SECRET": "abcdefgh"}
	opts := masker.DefaultOptions()
	opts.VisiblePrefix = 2
	opts.VisibleSuffix = 2

	results := masker.Mask(s, env, opts)
	if len(results) != 1 {
		t.Fatal("expected 1 result")
	}
	got := results[0].Masked
	if got[:2] != "ab" {
		t.Errorf("expected prefix 'ab', got %q", got[:2])
	}
	if got[len(got)-2:] != "gh" {
		t.Errorf("expected suffix 'gh', got %q", got[len(got)-2:])
	}
}

func TestMask_EmptyValueUnchanged(t *testing.T) {
	s := makeSchema(ev("EMPTY", true))
	env := map[string]string{"EMPTY": ""}
	opts := masker.DefaultOptions()

	results := masker.Mask(s, env, opts)
	if results[0].Masked != "" {
		t.Errorf("expected empty string, got %q", results[0].Masked)
	}
}

func TestMask_MissingEnvKeySkipped(t *testing.T) {
	s := makeSchema(ev("MISSING", true))
	env := map[string]string{}
	opts := masker.DefaultOptions()

	results := masker.Mask(s, env, opts)
	if len(results) != 0 {
		t.Errorf("expected no results for missing key, got %d", len(results))
	}
}

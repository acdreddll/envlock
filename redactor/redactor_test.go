package redactor_test

import (
	"testing"

	"github.com/envlock/redactor"
	"github.com/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key string, sensitive bool, def string) schema.EnvVar {
	return schema.EnvVar{Key: key, Sensitive: sensitive, Default: def}
}

func TestRedact_NonSensitivePassThrough(t *testing.T) {
	s := makeSchema(ev("APP_ENV", false, ""))
	env := map[string]string{"APP_ENV": "production"}

	results := redactor.Redact(s, env, redactor.Options{})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Value != "production" || results[0].Redacted {
		t.Errorf("expected plain value, got %+v", results[0])
	}
}

func TestRedact_SensitiveIsMasked(t *testing.T) {
	s := makeSchema(ev("DB_PASSWORD", true, ""))
	env := map[string]string{"DB_PASSWORD": "s3cr3t"}

	results := redactor.Redact(s, env, redactor.Options{})
	if !results[0].Redacted {
		t.Error("expected sensitive value to be redacted")
	}
	if results[0].Value != "***REDACTED***" {
		t.Errorf("unexpected mask: %s", results[0].Value)
	}
}

func TestRedact_CustomMask(t *testing.T) {
	s := makeSchema(ev("API_KEY", true, ""))
	env := map[string]string{"API_KEY": "abc123"}

	results := redactor.Redact(s, env, redactor.Options{Mask: "[hidden]"})
	if results[0].Value != "[hidden]" {
		t.Errorf("expected custom mask, got %s", results[0].Value)
	}
}

func TestRedact_RedactAll(t *testing.T) {
	s := makeSchema(ev("APP_ENV", false, ""), ev("DB_PASSWORD", true, ""))
	env := map[string]string{"APP_ENV": "staging", "DB_PASSWORD": "pass"}

	results := redactor.Redact(s, env, redactor.Options{RedactAll: true})
	for _, r := range results {
		if !r.Redacted {
			t.Errorf("expected %s to be redacted with RedactAll", r.Key)
		}
	}
}

func TestRedact_FallsBackToDefault(t *testing.T) {
	s := makeSchema(ev("LOG_LEVEL", false, "info"))
	env := map[string]string{}

	results := redactor.Redact(s, env, redactor.Options{})
	if results[0].Value != "info" {
		t.Errorf("expected default value 'info', got %s", results[0].Value)
	}
}

func TestRedact_EmptySensitiveNotMasked(t *testing.T) {
	s := makeSchema(ev("SECRET", true, ""))
	env := map[string]string{}

	results := redactor.Redact(s, env, redactor.Options{})
	if results[0].Redacted {
		t.Error("empty sensitive value should not be marked as redacted")
	}
}

func TestIndexByKey(t *testing.T) {
	s := makeSchema(ev("A", false, ""), ev("B", true, ""))
	env := map[string]string{"A": "1", "B": "secret"}

	results := redactor.Redact(s, env, redactor.Options{})
	idx := redactor.IndexByKey(results)

	if idx["A"].Value != "1" {
		t.Errorf("unexpected value for A: %s", idx["A"].Value)
	}
	if idx["B"].Value != "***REDACTED***" {
		t.Errorf("unexpected value for B: %s", idx["B"].Value)
	}
}

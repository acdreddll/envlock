package cascader_test

import (
	"testing"

	"github.com/envlock/cascader"
	"github.com/envlock/schema"
)

func makeSchema(vars ...schema.EnvVar) schema.Schema {
	return schema.Schema{Vars: vars}
}

func ev(key string, required bool, def string) schema.EnvVar {
	return schema.EnvVar{Key: key, Required: required, Default: def}
}

func TestCascade_FirstSourceWins(t *testing.T) {
	s := makeSchema(ev("DB_HOST", true, ""))
	sources := []cascader.Source{
		{Name: "override", Values: map[string]string{"DB_HOST": "override-host"}},
		{Name: "env", Values: map[string]string{"DB_HOST": "env-host"}},
	}
	results, errs := cascader.Cascade(s, sources)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if results[0].Value != "override-host" || results[0].Source != "override" {
		t.Errorf("expected override-host from override, got %q from %q", results[0].Value, results[0].Source)
	}
}

func TestCascade_FallsBackToLaterSource(t *testing.T) {
	s := makeSchema(ev("API_KEY", true, ""))
	sources := []cascader.Source{
		{Name: "override", Values: map[string]string{}},
		{Name: "env", Values: map[string]string{"API_KEY": "secret"}},
	}
	results, errs := cascader.Cascade(s, sources)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if results[0].Value != "secret" || results[0].Source != "env" {
		t.Errorf("expected secret from env, got %q from %q", results[0].Value, results[0].Source)
	}
}

func TestCascade_UsesDefault(t *testing.T) {
	s := makeSchema(ev("LOG_LEVEL", false, "info"))
	sources := []cascader.Source{
		{Name: "env", Values: map[string]string{}},
	}
	results, errs := cascader.Cascade(s, sources)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if results[0].Value != "info" || results[0].Source != "default" {
		t.Errorf("expected info from default, got %q from %q", results[0].Value, results[0].Source)
	}
}

func TestCascade_RequiredMissingReturnsError(t *testing.T) {
	s := makeSchema(ev("SECRET", true, ""))
	sources := []cascader.Source{
		{Name: "env", Values: map[string]string{}},
	}
	_, errs := cascader.Cascade(s, sources)
	if len(errs) == 0 {
		t.Fatal("expected an error for missing required variable")
	}
}

func TestCascade_OptionalMissingNoError(t *testing.T) {
	s := makeSchema(ev("OPTIONAL_VAR", false, ""))
	sources := []cascader.Source{
		{Name: "env", Values: map[string]string{}},
	}
	results, errs := cascader.Cascade(s, sources)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if results[0].Resolved {
		t.Error("expected unresolved for optional missing var")
	}
}

func TestToMap_OnlyResolved(t *testing.T) {
	results := []cascader.Result{
		{Key: "A", Value: "1", Source: "env", Resolved: true},
		{Key: "B", Value: "", Source: "", Resolved: false},
	}
	m := cascader.ToMap(results)
	if _, ok := m["B"]; ok {
		t.Error("unresolved key B should not appear in map")
	}
	if m["A"] != "1" {
		t.Errorf("expected A=1, got %q", m["A"])
	}
}
